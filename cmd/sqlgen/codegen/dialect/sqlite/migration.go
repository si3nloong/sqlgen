package sqlite

// func (s *sqliteDriver) Migrate(ctx context.Context, w io.Writer, t *compiler.Table) error {
// 	// TODO: Need to add up and down migration
// 	fmt.Fprint(w, "CREATE TABLE "+s.QuoteIdentifier(t.Name)+" (\n")
// 	if n := len(t.Columns); n > 0 {
// 		column := t.Columns[0]
// 		fmt.Fprint(w, column.Name, column.GoType().String())
// 		for i := 1; i < n; i++ {
// 			column = t.Columns[i]
// 			if column.Readonly {
// 				continue
// 			}
// 			if column.IsNullable() {
// 				fmt.Fprint(w, ",\n"+column.Name)
// 			} else {
// 				fmt.Fprint(w, ",\n"+column.Name+"VARCHAR(255) NOT NULL")
// 			}
// 		}
// 	}
// 	if pk, ok := t.PK(); ok {
// 		switch v := pk.(type) {
// 		case *compiler.PrimaryKey:
// 			log.Println(v)
// 		case *compiler.CompositePrimaryKey:
// 			log.Println(v)
// 		}
// 		fmt.Fprintf(w, "PRIMARY KEY ();")
// 	}
// 	fmt.Fprintf(w, ");")
// 	return nil
// }
