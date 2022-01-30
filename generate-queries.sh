export POSTGRES_PASSWORD=$(kubectl get secret --namespace default psql-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)
kubectl port-forward --namespace default svc/psql-postgresql 5432:5432 &
sleep 2
PGPASSWORD="$POSTGRES_PASSWORD" pg_dump --host 127.0.0.1 -U postgres -d postgres -p 5432 -s > schema.sql
rm -Rf sqlc
sqlc generate
kill $(jobs -p)