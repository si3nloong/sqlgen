# 常见问题解答（FAQ）

### 我如何声明我的模型？

您可以参考[这个指南](./MODELS.md)。

### 我不想在我的表上有主键，我可以吗？

当然可以。如果在结构标签中未指定`pk`、`primary_key`或`auto_increment`值，那么将不会生成主键方法。

### `sqlgen`是否支持`UUID`主键？

当然支持。`sqlgen`甚至支持有序 UUID。有关实现细节，您可以参考[这里](./UUID.md)。

### 我的结构实现了一些 sequel 接口，例如`sequel.Tabler`、`sequel.Columner`等，它们会被`sqlgen`覆盖吗？

不会，`sqlgen`将尊重结构方法的手动实现。
