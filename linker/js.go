package linker

import (
	"github.com/nar-lang/nar-common/bytecode"
	"github.com/nar-lang/nar-common/logger"
	"github.com/nar-lang/nar-compiler/locator"
)

func NewJsLinker(cleanup bool, cacheDir string) Linker {
	return &JsLinker{cleanup: cleanup, cacheDir: cacheDir}
}

type JsLinker struct {
	cleanup  bool
	cacheDir string
}

//pack := flag.String("pack", "", "command to pack resulted executable.\n"+ "  examples\n"+"  js: `webpack-cli --entry build/index.source.js -o ./build`")

func (l *JsLinker) Link(log *logger.LogWriter, binary *bytecode.Binary, lc locator.Locator, debug bool) error {
	//runtimePath, err := cacheRuntime(log, l.cacheDir)
	//if err != nil {
	//	return err
	//}
	//
	//indexJs := strings.Builder{}
	//var nativeNames []string
	//indexJs.WriteString(fmt.Sprintf("import NarRuntime from '%s'\n", runtimePath))
	//
	//for _, pkg := range packages {
	//	_, err := os.Stat(path.Join(pkg.Dir, nativeIndexPath))
	//	if os.IsNotExist(err) {
	//		continue
	//	}
	//	if err != nil {
	//		return err
	//	}
	//
	//	jsName := dotToUnderscore(pkg.Package.Name)
	//	nativeNames = append(nativeNames, jsName)
	//	indexJs.WriteString(fmt.Sprintf("import %s from '%s'\n", jsName, filepath.Join(pkg.Dir, nativeIndexPath)))
	//}
	//
	//indexJs.WriteString("const req = new XMLHttpRequest();\n")
	//indexJs.WriteString("req.open('GET', 'program.binar', true);\n")
	//indexJs.WriteString("req.responseType = 'arraybuffer';\n")
	//indexJs.WriteString("req.onload = function(e) {\n")
	//indexJs.WriteString("    const runtime = new NarRuntime(e.target.response);\n")
	//
	//for _, name := range nativeNames {
	//	indexJs.WriteString(fmt.Sprintf("    %s(runtime);\n", name))
	//}
	//
	//if main != "" {
	//	indexJs.WriteString(fmt.Sprintf("    runtime.execute('%s');\n", main))
	//}
	//
	//indexJs.WriteString("};\nreq.send(null);\n")
	//
	//l.artifactPath = filepath.Join(out, "main.source.js")
	//
	//err = os.WriteFile(l.artifactPath, []byte(indexJs.String()), 0640)
	//if err != nil {
	//	return logger.NewSystemError(err)
	//}
	//
	//htmlPath := filepath.Join(out, "index.html")
	//err = os.WriteFile(htmlPath, indexHtml, 0640)
	//if err != nil {
	//	return logger.NewSystemError(err)
	//}
	//
	//logger.Trace("linked successfully")
	return nil
}

//func dotToUnderscore(s ast.PackageIdentifier) string {
//	return strings.ReplaceAll(string(s), ".", "_")
//}
//
//func cacheRuntime(cacheDir string, upgrade bool, logger *logger.LogWriter) (string, error) {
//	runtimeDir, err := filepath.Abs(filepath.Join(cacheDir, "runtime-js"))
//	if err != nil {
//		return "", err
//	}
//	runtimePath := filepath.Join(runtimeDir, "index.js")
//	_, err = os.Stat(runtimePath)
//	if err != nil && !os.IsNotExist(err) {
//		return "", err
//	}
//	loaded := err == nil
//
//	if !loaded {
//		logger.Trace(fmt.Sprintf("cloning runtime `%s`\n", runtimeRepositoryUrl))
//		w := bytes.NewBufferString("")
//		_, err := git.PlainClone(runtimeDir, false, &git.CloneOptions{
//			URL:      runtimeRepositoryUrl,
//			Progress: w,
//		})
//		logger.Trace(w.String())
//		if err != nil {
//			return "", err
//		}
//		logger.Trace(fmt.Sprintf("downloaded runtime `%s`\n", runtimeRepositoryUrl))
//	} else if upgrade {
//		r, err := git.PlainOpen(runtimeDir)
//		if err == nil {
//			logger.Trace(fmt.Sprintf("upgrading runtime `%s` ", runtimeRepositoryUrl))
//			worktree, err := r.Worktree()
//			w := bytes.NewBufferString("")
//			if err != nil {
//				logger.Warn(err)
//			} else {
//				err = worktree.Pull(&git.PullOptions{
//					Progress: w,
//				})
//				logger.Trace(w.String())
//				if err != nil {
//					logger.Warn(err)
//				} else {
//					logger.Info(fmt.Sprintf("runtime upgraded `%s`", runtimeRepositoryUrl))
//				}
//			}
//		}
//	}
//	return runtimePath, nil
//}
//
//var indexHtml = []byte(
//	"<!DOCTYPE html>\n" +
//		"<html lang=\"en\">\n" +
//		"<head>\n" +
//		"    <meta charset=\"UTF-8\">\n" +
//		"    <title>Nar test</title>\n" +
//		"    <script src=\"main.js\" type=\"module\"></script>\n" +
//		"</head>\n" +
//		"<body></body>\n" +
//		"</html>\n")
