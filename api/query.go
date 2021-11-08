package api

import (
	"fmt"
	"time"
)

type QueryBuilder struct {
	Query string
}

func NewQueryBuilder() QueryBuilder {
	return QueryBuilder{Query: ""}
}

func addQuery(qb *QueryBuilder, query string, join string) {
	if qb.Query == "" {
		qb.Query = query
	} else {
		qb.Query = fmt.Sprintf("%s %s %s", qb.Query, join, query)
	}
}

// And joins the query with an and operator
func (qb *QueryBuilder) And(query string) *QueryBuilder {
	addQuery(qb, query, "AND")
	return qb
}

// Or joins the query with an or operator
func (qb *QueryBuilder) Or(query string) *QueryBuilder {
	addQuery(qb, query, "OR")
	return qb
}

func (qb *QueryBuilder) GetQuery() string {
	return qb.Query
}

func opGroup(queries []string, join string) string {
	query := ""
	for _, q := range queries {
		if query == "" {
			query = "(" + q
		} else {
			query = fmt.Sprintf("%s %s %s", query, join, q)
		}
	}
	return query + ")"
}

func AndGroup(queries ...string) string {
	return opGroup(queries, "AND")
}

func OrGroup(queries ...string) string {
	return opGroup(queries, "OR")
}

func IntersectGeoPoint(geoPoint GeoPoint) string {
	newQuery := "footprint:\"Intersects(%f,%f)\""
	return fmt.Sprintf(newQuery, geoPoint.Latitude, geoPoint.Longitude)
}

func PlatformName(platformName string) string {
	newQuery := "platformname:%s"
	return fmt.Sprintf(newQuery, platformName)
}

func SensingRange(start time.Time, end time.Time) string {
	newQuery := "beginposition:[%s TO %s]"
	return fmt.Sprintf(newQuery, start.Format(time.RFC3339), end.Format(time.RFC3339))
}

func ProductType(productType string) string {
	newQuery := "producttype:%s"
	return fmt.Sprintf(newQuery, productType)
}
