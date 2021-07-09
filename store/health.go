package store

// // type check is a db health check struct
// type check struct {
// 	config config.DB
// 	*sql.DB
// }

// // NewHealth creates new check Instance
// func NewHealth(c config.DB, db *sql.DB) health.Health {
// 	return &check{c, db}
// }

// // Full is checking connection to db by calling Ping() method.
// // it will return health.Response struct
// func (c *check) Full() health.Response {
// 	return health.Response{
// 		Name:    c.config.Name,
// 		Time:    time.Now().UTC(),
// 		Healthy: c.Ping() == nil,
// 	}
// }
