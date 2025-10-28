# ESC键导航功能实现说明

## 概述
实现了ESC键导航功能，可以从所有页面返回到主页（Home），不会退出应用程序。

## 修改内容

### 1. 修改了 `internal/ui/app.go`

#### globalInputHandler 函数
添加了ESC键处理逻辑：

```go
// Handle ESC key - return to home page (except when already on home page)
if event.Key() == utils.CloseDialogKey.Key {
    if a.currentPageIdx != homePageIndex {
        a.switchToPage(homePageIndex)
        return nil
    }
    // If already on home page, do nothing (don't exit)
    return nil
}
```

**主要特性：**
- 首先检查是否有对话框获得焦点（对话框自行处理ESC关闭）
- 如果在任何非主页的页面，ESC键返回主页
- 如果已经在主页，ESC键不做任何操作（不退出应用）
- 对话框的ESC处理优先于页面导航

#### updateHelpBar 函数
更新了帮助栏，在除主页外的所有页面显示"ESC Home"：

```go
case terminalPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Enter[-] Execute"
case databasePageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Ctrl+N[-] Connect | ..."
case toolsPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Space[-] Toggle | ..."
case settingsPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Space[-] Toggle | ..."
case systemPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]↑/↓[-] Scroll"
```

## 导航流程

### ESC键优先级顺序：
1. **对话框打开时**：ESC关闭对话框（停留在当前页面）
2. **非主页页面**：ESC返回主页
3. **主页**：ESC不做任何操作（应用保持打开）

### 示例流程：
```
主页 (F1) → 数据库页面 (F4) → [ESC] → 主页 (F1)
主页 (F1) → 工具页面 (F6) → [ESC] → 主页 (F1)
数据库页面 → Ctrl+N (打开对话框) → [ESC] (关闭对话框) → [ESC] → 主页
```

## 用户体验改进

### 修改前：
- ESC键在各页面没有一致的行为
- 用户可能意外按ESC导致退出
- 没有快速返回主页的明确方式

### 修改后：
- ✅ 所有页面ESC行为一致
- ✅ 安全导航 - ESC永远不会退出应用
- ✅ 直观 - ESC返回主页面
- ✅ 对话框感知 - ESC优先关闭对话框
- ✅ 视觉反馈 - 帮助栏显示"ESC Home"

## 键盘快捷键汇总

| 按键 | 功能 | 说明 |
|-----|------|------|
| **ESC** | 返回主页 | 在主页无效；优先关闭对话框 |
| **F1** | 进入主页 | 直接导航 |
| **F2** | 进入终端页面 | 直接导航 |
| **F4** | 进入数据库页面 | 直接导航 |
| **F6** | 进入工具页面 | 直接导航 |
| **F7** | 进入设置页面 | 直接导航 |
| **F8** | 进入系统信息页面 | 直接导航 |
| **Tab** | 循环切换页面 | 顺序导航 |
| **q** | 退出应用 | 完全退出 |

## 技术细节

### 实现模式：
ESC键处理器放在全局输入处理链的早期位置：
1. 检查对话框焦点 → 如果对话框打开，将事件传递给对话框
2. 检查ESC键 → 处理导航
3. 处理其他全局按键（q、Tab、功能键）

这确保了正确的事件冒泡，避免冲突。

### 代码质量：
- ✅ 不破坏现有功能
- ✅ 保持现有对话框行为
- ✅ 遵循项目代码风格和模式
- ✅ 使用现有常量（`utils.CloseDialogKey`）
- ✅ 成功编译，无错误
- ✅ 所有页面测试通过

## 测试

详细的测试用例和验证步骤请参见 [TEST_ESC_NAVIGATION.md](../TEST_ESC_NAVIGATION.md)

## 兼容性

- **Windows**：已测试，正常工作
- **Linux**：应该正常工作（使用tcell库）
- **macOS**：应该正常工作（使用tcell库）

实现使用跨平台的`tcell`库，在所有平台上都能一致地处理ESC键（KeyEscape）。

## 使用建议

1. **日常使用**：在任何页面按ESC快速返回主页
2. **对话框操作**：按ESC关闭对话框，再按ESC返回主页
3. **退出应用**：使用'q'键安全退出，不要担心误按ESC退出

## 文件修改清单

- `internal/ui/app.go` - 添加ESC键处理逻辑和更新帮助栏
- `docs/ESC_KEY_NAVIGATION.md` - 英文文档（新增）
- `docs/ESC键导航说明.md` - 中文文档（本文件，新增）
- `TEST_ESC_NAVIGATION.md` - 测试指南（新增）

## 构建和运行

```bash
# 进入项目目录
cd gocmder

# 编译
go build -o gocmder.exe .

# 运行
.\gocmder.exe
```

## 验证方法

1. 启动应用，当前在主页（F1）
2. 按F6进入工具页面
3. 按ESC，应该返回主页
4. 再按ESC，应该停留在主页（不退出）
5. 按'q'退出应用

完成！所有页面的ESC键都能正确返回主页，主页按ESC不会退出应用。

