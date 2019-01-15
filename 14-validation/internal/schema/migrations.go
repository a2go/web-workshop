package schema

import "github.com/GuiaBolso/darwin"

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Including the queries directly in this file has the same pros/cons mentioned
// in seeds.go

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script: `
CREATE TABLE products (
	product_id UUID,
	name       VARCHAR(255),
	cost       INT,
	quantity   INT,

	PRIMARY KEY (product_id)
);`,
	},
	{
		Version:     2,
		Description: "Add sales",
		Script: `
CREATE TABLE sales (
	sale_id    UUID,
	product_id UUID,
	amount     INT,
	paid       INT,

	PRIMARY KEY (sale_id),
	FOREIGN KEY (product_id) REFERENCES products(product_id)
);`,
	},
}
