package gsd

import (
	"fmt"
	"strings"
)

/********** mssqlBuilder **********/

// SQL Server 2012+ builder
type mssqlBuilder struct {
}

// BuildInsert build query string and parameters for insert action
func (this *mssqlBuilder) BuildInsert(ctx *buildContext, info *insertInfo) error {
	ctx.AppendSqlF("INSERT INTO [%s](", info.table)

	first := true
	for k, v := range info.values {
		if first {
			first = false
		} else {
			ctx.AppendSql(",")
		}
		ctx.AppendSqlF("[%s]", k)
		ctx.AddParam(v)
	}

	ctx.AppendSql(") VALUES(?", strings.Repeat(",?", len(info.values)-1), ")")

	return nil
}

// BuildUpdate build query string and parameters for update action
func (this *mssqlBuilder) BuildUpdate(ctx *buildContext, info *updateInfo) error {
	ctx.AppendSqlF("UPDATE [%s] SET", info.table)

	first := true
	for k, v := range info.values {
		if first {
			first = false
		} else {
			ctx.AppendSql(",")
		}

		switch v.ut {
		case UPDATE_INC:
			ctx.AppendSqlF(" [%s]=[%s]+?", k, k)
			ctx.AddParam(v.val)
		case UPDATE_XP:
			ctx.AppendSqlF(" [%s]=%s", k, v.val)
		default:
			ctx.AppendSqlF(" [%s]=?", k)
			ctx.AddParam(v.val)
		}
	}

	if info.where != nil {
		ctx.AppendSql(" WHERE ")
		this.BuildFilters(ctx, info.where)
	}

	return nil
}

// BuildDelete build query string and parameters for delete action
func (this *mssqlBuilder) BuildDelete(ctx *buildContext, info *deleteInfo) error {
	ctx.AppendSqlF("DELETE FROM [%s]", info.table)

	if info.where != nil {
		ctx.AppendSql(" WHERE ")
		this.BuildFilters(ctx, info.where)
	}

	return nil
}

// BuildSelect build query string and parameters for select action
func (this *mssqlBuilder) BuildSelect(ctx *buildContext, info *selectInfo) error {
	ctx.AppendSql("SELECT ")

	if info.distinct {
		ctx.AppendSql("DISTINCT ")
	}

	// SELECT
	for i, c := range info.columns {
		if i > 0 {
			ctx.AppendSql(",")
		}

		switch v := c.(type) {
		case *normalColumn:
			if v.table != nil {
				ctx.AppendSqlF("[%s].", v.table.Prefix())
			}
			ctx.AppendSqlF("[%s]", v.column)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		case *exprColumn:
			ctx.AppendSql(v.expr)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		}
	}

	// FROM
	ctx.AppendSqlF(" FROM [%s]", info.table.Name())
	if info.table.Alias() != "" {
		ctx.AppendSql(" AS ", info.table.Alias())
	}

	// JOIN
	for _, j := range info.joins {
		ctx.AppendSqlF(" %s [%s]", j.jt, j.t.Name())
		if j.t.Alias() != "" {
			ctx.AppendSql(" AS ", j.t.Alias())
		}
		ctx.AppendSql(" ON ")
		this.BuildFilters(ctx, j.on)
	}

	if info.where != nil {
		ctx.AppendSql(" WHERE ")
		this.BuildFilters(ctx, info.where)
	}

	// GROUP BY
	if len(info.groups) > 0 {
		ctx.AppendSql(" GROUP BY ")
		for i, g := range info.groups {
			if i > 0 {
				ctx.AppendSql(",")
			}
			for j, col := range g.columns {
				if j > 0 {
					ctx.AppendSql(",")
				}
				if g.table != nil {
					ctx.AppendSqlF("[%s].", g.table.Prefix())
				}
				ctx.AppendSqlF("[%s]", col)
			}
		}

		if info.having != nil {
			ctx.AppendSql(" HAVING ")
			this.BuildFilters(ctx, info.having)
		}
	}

	// ORDER BY
	if len(info.orders) > 0 {
		ctx.AppendSql(" ORDER BY ")
		for i, order := range info.orders {
			if i > 0 {
				ctx.AppendSql(",")
			}
			for j, col := range order.columns {
				if j > 0 {
					ctx.AppendSql(",")
				}
				if order.table != nil {
					ctx.AppendSqlF("[%s].", order.table.Prefix())
				}
				ctx.AppendSqlF("[%s]", col)
			}
			ctx.AppendSqlF(" %s", order.st)
		}
	}

	// LIMIT
	if info.skip != 0 || info.take != 0 {
		ctx.AppendSqlF(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", info.skip, info.take)
	}

	return nil
}

func (this *mssqlBuilder) BuildFilters(ctx *buildContext, filters Filters) error {
	switch v := filters.(type) {
	case *basicFilters:
		for i, f := range v.items {
			if i > 0 {
				ctx.AppendSql(" AND ")
			}
			err := this.BuildFilter(ctx, f)
			if err != nil {
				return err
			}
		}
	case *notFilters:
		ctx.AppendSql("NOT(")
		err := this.BuildFilters(ctx, v.inner)
		if err != nil {
			return err
		}
		ctx.AppendSql(")")
	case *andFilters:
		ctx.AppendSql("(")
		this.BuildFilters(ctx, v.left)
		ctx.AppendSql(") AND (")
		this.BuildFilters(ctx, v.right)
		ctx.AppendSql(")")
	case *orFilters:
		ctx.AppendSql("(")
		err := this.BuildFilters(ctx, v.left)
		if err != nil {
			return err
		}
		ctx.AppendSql(") OR (")
		err = this.BuildFilters(ctx, v.right)
		if err != nil {
			return err
		}
		ctx.AppendSql(")")
	}

	return nil
}

func (this *mssqlBuilder) BuildFilter(ctx *buildContext, filter interface{}) (err error) {
	switch f := filter.(type) {
	case *oneColumnFilter:
		err = this.BuildOneColumnFilter(ctx, f)
	case *twoColumnFilter:
		err = this.BuildTwoColumnFilter(ctx, f)
	case *exprFilter:
		ctx.AppendSql(f.expr)
	default:
		err = fmt.Errorf("invalid filter: %v", filter)
	}
	return
}

func (this *mssqlBuilder) BuildOneColumnFilter(ctx *buildContext, f *oneColumnFilter) error {
	if f.table != nil {
		ctx.AppendSqlF("[%s].", f.table.Prefix())
	}

	switch f.ft {
	case FILTER_NE:
		if f.value == nil {
			ctx.AppendSqlF("[%s] IS NOT NULL", f.column)
		} else {
			ctx.AppendSqlF("[%s]<>?", f.column)
			ctx.AddParam(f.value)
		}
	case FILTER_LT:
		ctx.AppendSqlF("[%s]<?", f.column)
		ctx.AddParam(f.value)
	case FILTER_GT:
		ctx.AppendSqlF("[%s]>?", f.column)
		ctx.AddParam(f.value)
	case FILTER_LTE:
		ctx.AppendSqlF("[%s]<=?", f.column)
		ctx.AddParam(f.value)
	case FILTER_GTE:
		ctx.AppendSqlF("[%s]>=?", f.column)
		ctx.AddParam(f.value)
	case FILTER_IN:
		ctx.AppendSqlF("[%s] IN(%s)", f.column, f.value)
	case FILTER_LK:
		ctx.AppendSqlF("[%s] LIKE '%' + ? + '%'", f.column)
		ctx.AddParam(f.value)
	default:
		if f.value == nil {
			ctx.AppendSqlF("[%s] IS NULL", f.column)
		} else {
			ctx.AppendSqlF("[%s]=?", f.column)
			ctx.AddParam(f.value)
		}
	}

	return nil
}

func (this *mssqlBuilder) BuildTwoColumnFilter(ctx *buildContext, f *twoColumnFilter) error {
	switch f.ft {
	case FILTER_NE:
		ctx.AppendSqlF("[%s].[%s]<>[%s].[%s]", f.table1.Prefix(), f.column1, f.table2.Prefix(), f.column2)
	case FILTER_LT:
		ctx.AppendSqlF("[%s].[%s]<[%s].[%s]", f.table1.Prefix(), f.column1, f.table2.Prefix(), f.column2)
	case FILTER_GT:
		ctx.AppendSqlF("[%s].[%s]>[%s].[%s]", f.table1.Prefix(), f.column1, f.table2.Prefix(), f.column2)
	case FILTER_LTE:
		ctx.AppendSqlF("[%s].[%s]<=[%s].[%s]", f.table1.Prefix(), f.column1, f.table2.Prefix(), f.column2)
	case FILTER_GTE:
		ctx.AppendSqlF("[%s].[%s]>=[%s].[%s]", f.table1.Prefix(), f.column1, f.table2.Prefix(), f.column2)
	case FILTER_IN:
		return fmt.Errorf("invalid filterType: IN")
	case FILTER_LK:
		return fmt.Errorf("invalid filterType: LK")
	default:
		ctx.AppendSqlF("[%s].[%s]=[%s].[%s]", f.table1.Prefix(), f.column1, f.table2.Prefix(), f.column2)
	}

	return nil
}

/********** mssqlBuilder **********/

// SQL Server 2005+ builder
type mssql2005Builder struct {
	mssqlBuilder
}

// BuildSelect build query string and parameters for select action
func (this *mssql2005Builder) BuildSelect(ctx *buildContext, info *selectInfo) error {
	if info.skip == 0 {
		return this.buildSelectNoPage(ctx, info)
	} else {
		return this.buildSelectPage(ctx, info)
	}
}

// BuildSelect build query string and parameters for select action
func (this *mssql2005Builder) buildSelectNoPage(ctx *buildContext, info *selectInfo) error {
	ctx.AppendSql("SELECT ")

	if info.distinct {
		ctx.AppendSql("DISTINCT ")
	}

	// LIMIT
	if info.take > 0 {
		ctx.AppendSqlF("TOP %d ", info.take)
	}

	// SELECT
	for i, c := range info.columns {
		if i > 0 {
			ctx.AppendSql(",")
		}

		switch v := c.(type) {
		case *normalColumn:
			if v.table != nil {
				ctx.AppendSqlF("[%s].", v.table.Prefix())
			}
			ctx.AppendSqlF("[%s]", v.column)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		case *exprColumn:
			ctx.AppendSql(v.expr)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		}
	}

	// FROM
	ctx.AppendSqlF(" FROM [%s]", info.table.Name())
	if info.table.Alias() != "" {
		ctx.AppendSql(" AS ", info.table.Alias())
	}

	// JOIN
	for _, j := range info.joins {
		ctx.AppendSqlF(" %s [%s]", j.jt, j.t.Name())
		if j.t.Alias() != "" {
			ctx.AppendSql(" AS ", j.t.Alias())
		}
		ctx.AppendSql(" ON ")
		this.BuildFilters(ctx, j.on)
	}

	if info.where != nil {
		ctx.AppendSql(" WHERE ")
		this.BuildFilters(ctx, info.where)
	}

	// GROUP BY
	if len(info.groups) > 0 {
		ctx.AppendSql(" GROUP BY ")
		for i, g := range info.groups {
			if i > 0 {
				ctx.AppendSql(",")
			}
			for j, col := range g.columns {
				if j > 0 {
					ctx.AppendSql(",")
				}
				if g.table != nil {
					ctx.AppendSqlF("[%s].", g.table.Prefix())
				}
				ctx.AppendSqlF("[%s]", col)
			}
		}

		if info.having != nil {
			ctx.AppendSql(" HAVING ")
			this.BuildFilters(ctx, info.having)
		}
	}

	// ORDER BY
	if len(info.orders) > 0 {
		ctx.AppendSql(" ORDER BY ")
		for i, order := range info.orders {
			if i > 0 {
				ctx.AppendSql(",")
			}
			for j, col := range order.columns {
				if j > 0 {
					ctx.AppendSql(",")
				}
				if order.table != nil {
					ctx.AppendSqlF("[%s].", order.table.Prefix())
				}
				ctx.AppendSqlF("[%s]", col)
			}
			ctx.AppendSqlF(" %s", order.st)
		}
	}

	return nil
}

// BuildSelect build query string and parameters for select action
func (this *mssql2005Builder) buildSelectPage(ctx *buildContext, info *selectInfo) error {
	ctx.AppendSql("SELECT ")

	// SELECT
	for i, c := range info.columns {
		if i > 0 {
			ctx.AppendSql(",")
		}

		switch v := c.(type) {
		case *normalColumn:
			if v.table != nil {
				ctx.AppendSqlF("[%s].", v.table.Prefix())
			}
			ctx.AppendSqlF("[%s]", v.column)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		case *exprColumn:
			ctx.AppendSql(v.expr)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		}
	}

	ctx.AppendSql(" FROM (SELECT ")

	if info.distinct {
		ctx.AppendSql("DISTINCT ")
	}

	// SELECT
	for i, c := range info.columns {
		if i > 0 {
			ctx.AppendSql(",")
		}

		switch v := c.(type) {
		case *normalColumn:
			if v.table != nil {
				ctx.AppendSqlF("[%s].", v.table.Prefix())
			}
			ctx.AppendSqlF("[%s]", v.column)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		case *exprColumn:
			ctx.AppendSql(v.expr)
			if v.alias != "" {
				ctx.AppendSql(" AS ", v.alias)
			}
		}
	}

	// ORDER BY
	ctx.AppendSql(",ROW_NUMBER() OVER(")
	if len(info.orders) > 0 {
		ctx.AppendSql("ORDER BY ")
		for i, order := range info.orders {
			if i > 0 {
				ctx.AppendSql(",")
			}
			for j, col := range order.columns {
				if j > 0 {
					ctx.AppendSql(",")
				}
				if order.table != nil {
					ctx.AppendSqlF("[%s].", order.table.Prefix())
				}
				ctx.AppendSqlF("[%s]", col)
			}
			ctx.AppendSqlF(" %s", order.st)
		}
	}
	ctx.AppendSql(") AS _N")

	// FROM
	ctx.AppendSqlF(" FROM [%s]", info.table.Name())
	if info.table.Alias() != "" {
		ctx.AppendSql(" AS ", info.table.Alias())
	}

	// JOIN
	for _, j := range info.joins {
		ctx.AppendSqlF(" %s [%s]", j.jt, j.t.Name())
		if j.t.Alias() != "" {
			ctx.AppendSql(" AS ", j.t.Alias())
		}
		ctx.AppendSql(" ON ")
		this.BuildFilters(ctx, j.on)
	}

	if info.where != nil {
		ctx.AppendSql(" WHERE ")
		this.BuildFilters(ctx, info.where)
	}

	// GROUP BY
	if len(info.groups) > 0 {
		ctx.AppendSql(" GROUP BY ")
		for i, g := range info.groups {
			if i > 0 {
				ctx.AppendSql(",")
			}
			for j, col := range g.columns {
				if j > 0 {
					ctx.AppendSql(",")
				}
				if g.table != nil {
					ctx.AppendSqlF("[%s].", g.table.Prefix())
				}
				ctx.AppendSqlF("[%s]", col)
			}
		}

		if info.having != nil {
			ctx.AppendSql(" HAVING ")
			this.BuildFilters(ctx, info.having)
		}
	}

	ctx.AppendSqlF(") AS _T WHERE _N>%d AND _N<=%d", info.skip, info.skip+info.take)

	return nil
}
