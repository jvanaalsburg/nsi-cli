# NSI CLI

Command line interface for the NSI API.

### Configure and Build

Run the following command to build the application.

```bash
docker compose run --rm cli go build -o bin/nsi-cli
```

In order to user the application you will need a config file. Open the file `~/.config/nsi/nsi-cli.toml` and add an `auth` section with the NSI API URL root.

```toml
[api]
    url_root = "http://localhost:4141/nsiapi"
```

The above configuration will work with the default NSI API settings. If you are running the API on a different port, or if you want to connect to a different host, you will need to modify the URL value.

### Authentication

The `auth` command is used to log into an NSI API account. The API token will be stored in the `nsi-cli.toml` config file and will be automatically included in the authorization header when making requests.

```bash
nsi-cli auth login -email user@example.com
```

You can also specify the password when logging in. This should not be done for actual production accounts, but can be useful during development to quickly switch between accounts.

```bash
# Log in as the regular test user account.
nsi-cli auth login -email user@example.com -password secret

# Log in as the admin test user account.
nsi-cli auth login -email admin@example.com -password secret
```

There are a couple other `auth` actions that can be useful during development. The `status` command will display the account that is currently logged in.

```bash
nsi-cli auth status
```

The `token` action will return the saved auth token. This can be useful when making requests with another tool like `curl`.

```bash
curl -H "Authorization: Bearer $(nsi-cli auth token)" http://localhost:4141/users
```

### Development

If you want to run the application within the Docker container, you will need to create a config file similar to the one above. Create a `config/` directory, if it does not already exists, and add a `nsi-cli.dev.toml` file with the following contents:

```toml
[api]
    url_root = "http://api:4141/nsiapi"
```

**NOTE:** When running commands in a container, you must use `go run main.go` rather than `nsi-cli`. For example, the auth login command is called like this when using the `nsi-cli` executable:

```bash
nsi-cli auth login -email user@example.com
```

When running in the `cli` container, you would instead run the following:

```bash
docker compose run --rm cli go run main.go auth login -email user@example.com
```

### Examples

**Fetching Records**

```bash
# Get all platform users (must be logged in as an admin).
nsi-cli users list

# Get a specific user record.
nsi-cli users find -user-id "00000000-0000-0000-0000-000000000001"

# There are similar commands for getting group records.
nsi-cli groups list
nsi-cli groups find -group-id "00000000-0000-0000-0000-000000000001"
```

**Creating User Accounts**

```bash
nsi-cli users create \
    -email "john.doe@mail.com" \
    -first-name "John" \
    -last-name "Doe" \
    -password "secret"

# A random password will be generated if the password flag is omitted.
nsi-cli users create \
    -email "john.doe@mail.com" \
    -first-name "John" \
    -last-name "Doe"
```

**Managing Group Members**

```bash
# Add a user to a group.
nsi-cli groups add-user \
    -group-id "00000000-0000-0000-0000-000000000001" \
    -user-id "00000000-0000-0000-0000-000000000002" \
    -role "user"

# Remove a user from a group.
nsi-cli groups remove-user \
    -group-id "00000000-0000-0000-0000-000000000001" \
    -user-id "00000000-0000-0000-0000-000000000002"
```

### Generating UUIDs

Test records are typically assigned sequential IDs to make it easier to remember and identify them. However, it can still be cumbersome to write out the UUIDs when making requests, either with `nsi-cli` or another tool like `curl`. The `mkuuid` tool can be used to generate a UUID from a simple numeric ID. For example, the command `mkuuid 1` would return `00000000-0000-0000-0000-000000000001`.

**Examples of Using `mkuuid`**

```bash
nsi-cli users find -user-id $(mkuuid 1)
nsi-cli groups find -group-id $(mkuuid 2)

nsi-cli groups remove-user -group-id $(mkuuid 1) -user-id $(mkuuid 2)
nsi-cli groups add-user -group-id $(mkuuid 1) -user-id $(mkuuid 2) -role "user"
```

Run the following command to build the `mkuuid` tool.

```bash
docker compose run --rm cli go build -o bin/mkuuid ./tools/mkuuid
```
