package test

//func TestWithPostgres(t *testing.T) {
//	ctx := context.Background()
//	dbName := "postgres"
//	dbUser := "user"
//	dbPassword := "password"
//
//	pgC, err := postgres.Run(ctx,
//		"postgres:16-alpine",
//		postgres.WithDatabase(dbName),
//		postgres.WithUsername(dbUser),
//		postgres.WithPassword(dbPassword),
//		postgres.BasicWaitStrategies(),
//	)
//	defer func() {
//		if err := testcontainers.TerminateContainer(pgC); err != nil {
//			fmt.Printf("failed to terminate container: %s", err)
//		}
//	}()
//
//	require.NoError(t, err)
//}
