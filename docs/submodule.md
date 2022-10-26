## 初始化

```bash
git clone git@github.com:XRSec/Go-Wechaty-Bot.git
git submodule update --init --recursive
```

## 更新本地代码

```bash
git submodule update --recursive --remote
```

## 开发者

```bash
cd Server/Plug
git submodule update --recursive --remote
git checkout main
git add .
git commit -am "update"
cd ../..
git add .
git commit -am "update"
git push
```