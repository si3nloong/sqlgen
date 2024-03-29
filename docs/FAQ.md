# Frequent Ask Questions (FAQs)

### How can I declare my model?

You may refer to [this guide](./MODELS.md).

### How do I rename my table name?

By default, `sqlgen` is using struct name as table name, but we also offer `sequel.Table` for user to rename their table name. For more details, you may refer to [here](./TABLE_NAME.md).

### I don't want primary key on my table, can I do that?

Absolutely, yes. If you didn't specific the `pk`, `primary_key` or `auto_increment` value in struct tag, the primary key methods will not be generated.

### Can `sqlgen` support `UUID` primary key?

Why not? `sqlgen` even support ordered uuid. For implementation details, you may refer to [here](./UUID.md).

### My struct implement some of the sequel interface, such as `sequel.Tabler`, `sequel.Columner` etc, will them get override by `sqlgen`?

No, `sqlgen` will respect the manual implementation of the struct methods.
