# Database Patterns

This section has been reorganized into focused guides for better maintainability and clarity.

## Guides

### [Database Connection](./database-connection.md)
**Setting up database and cache connections**

- Database configuration (MySQL, PostgreSQL, SQLite)
- Redis configuration (node and cluster)
- Config structure and service context integration
- Connection with cache setup

### [Model Generation](./model-generation.md)
**Generating models from SQL schemas**

- Model generation from local SQL files
- Model generation from remote datasource
- Multiple datasources support
- Generated field constants and methods

### [Condition Builder](./condition-builder.md)
**Building type-safe query conditions**

- Using generated field constants
- Comparison and pattern matching operators
- Pagination and sorting operators
- Dynamic condition building
- Chain API for fluent conditions

### [CRUD Operations](./crud-operations.md)
**Using generated CRUD methods**

- Insert operations (single and bulk)
- Query operations (find, page, count)
- Update and delete operations
- Complex aggregation queries
- Table sharding patterns

### [Best Practices](./best-practices.md)
**Database usage guidelines**

- Security practices
- Performance optimization
- Error handling patterns
- Common code patterns
- Anti-patterns to avoid

## Quick Start

1. **Configure your database** in `etc/etc.yaml` - See [Database Connection](./database-connection.md)
2. **Place SQL files** in `desc/sql/` directory - See [Model Generation](./model-generation.md)
3. **Run `jzero gen`** to generate models
4. **Use generated methods** with condition builder - See [CRUD Operations](./crud-operations.md)

## External Resources

- [modelx Documentation](https://docs.jzero.io/component/modelx) - Complete modelx API reference
