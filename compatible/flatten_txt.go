package compatible

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/utils"
)

type FlattenTxt struct{}

func (f *FlattenTxt) Decode(file *os.File) (*dir.Dir, error) {
	metaDir := &dir.Dir{DirName: utils.FileNameStrip(file.Name())} // 文件夹 名字 将文件名字作为 文件夹名字

	scanner := bufio.NewScanner(file) //逐行扫描 文件
	for scanner.Scan() {
		line := scanner.Text()

		// 115：// 开头的txt文档
		if strings.HasPrefix(line, flattenTxtPrefix) { // 文本对比 HasPrefix测试字符串s是否以前缀开头。
			line = line[len(flattenTxtPrefix):] // line 就等于 line 自己从 115：// 开始之后 的字符串
		}

		parts := strings.Split(line, flattenTxtSplit) // 按照 | 给 line分割
		if len(parts) < normalSplitLen {              // 如果 少于4个以上的 part  （3个 | 分出4个 part）
			//return metaDir, errors.New("Format Error ")
			fmt.Printf("本行字符串有误 %s ", line)
			continue

		} else if len(parts) == normalSplitLen { // 如果等于 4个 就是只有一级目录 （3个 | 分出4个 part）
			metaDir.Files = append(metaDir.Files, line) // 直接向根文件目录追加 文件

		} else {
			dirParts := parts[normalSplitLen:]         // 取出 文件目录 4个 | 以后的字符串
			treeNode := rebuildTree(metaDir, dirParts) // 将文件目录 和 目录的字符串 传入 封装树

			// 在上一步 将所有目录都录入进去后，在返回 回来进行文件追加
			//  犹豫最后一步返回的  metaDir 已经是 最底层的 最下一层的目录  expectDir = innerDir   rebuildTree(expectDir, dirpaths[1:]) 共同作用的结果
			// 所以  treeNode 的目录 已经处于最底层了。 在目录上追加文件就能直接追加到 正确的位置了
			treeNode.Files = append(treeNode.Files, strings.Join(parts[:normalSplitLen], flattenTxtSplit)) // strings.Join 是一个拼接 函数 只是将字符串给拼接起来
			// 现在 treeNode 处于 最后一层 ，文件追加成功了。如何将成功的文件给 追加给 metaDir 呢 treeNode 即将消亡
			// treeNode是子节点，有父节点指向它的
			// 所以父节点的根节点是metaDir
			// 所以最后只返回根节点metaDir就行了

		}
	}
	return metaDir, nil
}

func rebuildTree(metaDir *dir.Dir, dirpaths []string) *dir.Dir {
	// 当传入的 文件目录 是空的时候，就直接返回 原本的 目录结构
	if len(dirpaths) == 0 {
		return metaDir // 这是遍历完全部后的 返回 不然都返回  return rebuildTree(expectDir, dirpaths[1:])
	}

	found := false
	expectDir := &dir.Dir{}
	for _, innerDir := range metaDir.Dirs { // 遍历 已经存在的 文件目录
		if innerDir.DirName == dirpaths[0] { // 如果找到
			found = true         // 写入找到
			expectDir = innerDir // 这一步 至关重要，将 dir 目录 重置为当前已经存在的目录，这样才能保持树结构，同样是第二次循环  以及下次循环
		}
	}

	if !found { // 如果 没找到
		expectDir.DirName = dirpaths[0]                // 对第一级目录赋值 取 dirpaths的第一个
		metaDir.Dirs = append(metaDir.Dirs, expectDir) // 将 第一个目录给追加到 文件夹上
	}

	return rebuildTree(expectDir, dirpaths[1:]) // 重新 返回 下一次 计算 当然是取消第一个目录了 例如 有1 2 3 级目录 第二次就是 2 3 级，第三次 就是3  第4次就是 空的了 所以第一句 if len(dirpaths) == 0 { 才发挥作用

	// 这里是二叉树 节点法

	// 这里只是返回了 最末端的  expectDir 给 treeNode用，本身节点没有消亡。等待 最后计算完毕 将二叉树节点返回给上一个函数即可

	// 重点复习 二叉树 节点法
}
