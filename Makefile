# Include your custom make commands here.
generate:  gql_gen

gql_gen:
	@echo 'Running gqlgen'; \
	(cd transport/gql && go run github.com/99designs/gqlgen generate --config gqlgen.yml --verbose); \



