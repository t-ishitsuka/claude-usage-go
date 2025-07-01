# CLAUDE.md

このファイルは、このリポジトリのコードを扱う際のClaude Code (claude.ai/code) へのガイダンスを提供します。

## 最重要事項

**すべてのやり取りは日本語で行うこと。**

## プロジェクト概要

claude-usage-goは、ローカルのJSONLファイル（`~/.claude/projects/`）からClaude APIの使用状況を分析し、トークン使用量と推定コストを計算するCLIツールです。ccusageにインスパイアされたGo実装です。

## ビルドと開発コマンド

### Makeを使用（推奨）

```bash
# フォーマット、リント、テスト、ビルド
make

# 個別コマンド
make build         # アプリケーションをビルド
make test          # テストを実行
make test-coverage # カバレッジ付きでテストを実行
make format        # gofmtでコードをフォーマット
make lint          # リンターを実行
make clean         # ビルド成果物を削除
make run           # ビルドして実行
make help          # 利用可能なすべてのコマンドを表示
```

### 手動コマンド

```bash
# プロジェクトをビルド
go build -o claude-usage-go

# 依存関係をダウンロード
go mod tidy

# ツールを実行
./claude-usage-go daily
./claude-usage-go monthly
./claude-usage-go session

# 一般的な開発コマンド
go fmt ./...  # コードをフォーマット
go vet ./...  # 静的解析を実行

# テストコマンド
go test ./...                    # すべてのテストを実行
go test -v ./...                 # 詳細出力付きですべてのテストを実行
go test ./internal/calculator    # 特定パッケージのテストを実行
go test -run TestCalculateCost ./internal/calculator  # 特定のテストを実行
go test -cover ./...             # カバレッジ付きでテストを実行
go test -race ./...              # レース検出器付きでテストを実行
```

## アーキテクチャ

コードベースは関心事の明確な分離を持つクリーンアーキテクチャに従っています：

1. **CLIレイヤー** (`/cmd/`): ユーザーインタラクションを処理するCobraベースのコマンド

   - 各コマンド（daily、monthly、session）は同じパターンに従います：オプション解析 → データ取得 → フィルタリング → 集計 → 表示
   - グローバルフラグは`root.go`で定義され、パッケージレベル変数でアクセスされます

2. **データモデル** (`/internal/models/`): コアタイプと価格データ

   - `TokenUsage`はすべてのトークンタイプ（入力、出力、キャッシュ作成/読み取り）を追跡
   - モデル価格は`pricing.go`にハードコード - Claude価格が変更されたらここを更新
   - モデルIDのパターン：`claude-{version}-{model}-{date}` (例：`claude-opus-4-20250514`)

3. **JSONLパーサー** (`/internal/parser/`): ファイルの読み取りと解析を処理

   - JSONL形式：各行にネストされた構造のメッセージが含まれる
   - キー解析ロジック：`type: "assistant"`と`message.role: "assistant"`および`message.usage`データを探す
   - バッファサイズは大きなファイルを扱うため10MBに設定

4. **コスト計算機** (`/internal/calculator/`): 使用量を集計しコストを計算

   - すべての集計関数はソート済みの結果を返す
   - コスト計算：`tokens / 1_000_000 * price_per_1M`
   - 日次、月次、セッション集計をサポート

5. **表示** (`/internal/display/`): 出力をテーブルまたはJSONとしてフォーマット
   - テーブル表示はANSIカラーとUnicodeボックス文字を使用
   - ブレークダウンモードは各期間内のモデル別コストを表示

## 主要な実装詳細

- **JSONL構造**: メッセージはトップレベルではなく`message`フィールド内にネストされている
- **モデル命名**: 表示には`GetModelShortName()`を使用（例：完全なIDではなく「Opus 4」）
- **日付フィルタリング**: Goのtime.Parseを「20060102」形式（YYYYMMDD）で使用
- **モデルフィルタリング**: モデル名の大文字小文字を区別しない比較
- **外部APIなし**: すべてのデータはローカルのJSONLファイルから取得、ネットワーク呼び出しなし

## 新機能の追加

新しいClaudeモデルを追加する場合：

1. `internal/models/pricing.go`の`ModelPricing`マップに価格を追加
2. `GetModelShortName()`関数に短縮名マッピングを追加

JSONL解析を修正する場合：

1. 実際のJSONL構造を確認：`grep '"type":"message"' ~/.claude/projects/*/**.jsonl | head -1 | python3 -m json.tool`
2. `internal/parser/jsonl.go`の`JSONLEntry`と関連構造体を更新

## 現在の制限事項

- テストファイルが存在しない - 変更時にはテストの追加を検討
- 自動リンティングやCI/CDセットアップがない
- 価格データがハードコード（APIから取得しない）
- `~/.claude/projects/`ディレクトリからのみ読み取り
