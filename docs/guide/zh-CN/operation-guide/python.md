## python 示例

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
import requests

resp = requests.get("https://www.baidu.com/",
        proxies={
        "https": "http://localhost:8080",
        "http": "http://localhost:8080",
        },
        verify=False
   )
print(resp.text)
```