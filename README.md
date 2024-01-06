# go-pbt-demo
Goでプロパティベーステスト(PBT)を実践するためのデモ

## 書籍貸出システム

### システム要件

- 「まだシステムに登録されていない本を追加する」に期待されるのは「成功」
- 「すでにシステムに登録されている本を追加する」に期待されるのは「失敗」
- 「すでにシステムに登録されている本の在庫を 1 冊追加する」に期待されるのは「成功(すぐに在庫が 1 冊増える)」 
- 「まだシステムに登録されていない本の在庫を 1 冊追加する」に期待されるのは「エラー」 
- 「システムに登録されていて利用可能な在庫がある本を貸出する」に期待されるのは「在庫が 1 冊減る」 
- 「システムに登録されているが利用可能な在庫がない本を貸出する」に期待されるのは「貸出不能のエラー」 
- 「システムに登録されていない本を貸出する」に期待されるのは「書籍がないというエラー」 
- 「システムに登録されている本を返却する」に期待されるのは「在庫を戻す」 
- 「システムに登録されてない本を返却する」に期待されるのは「在庫がないというエラー」 
- 「システムに登録されていて利用可能な在庫が減っていない本を返却する」に期待されるのも「エラー」
- 「ISBN で本を検索する」に対し「その本がシステムに登録されている場合」に期待されるのは「成功」
- 「ISBN で本を検索する」に対し「その本がシステムに登録されていない場合」に期待されるのは「失敗」 
- 「著者名で本を検索する」に対し「著者名の一部または全体と一致する本が少なくとも 1 つ登録されている」に期待されるのは「成功」 
- 「書名で検索する」に対し「書名の一部または全部に一致する本が少なくとも 1 つ登録されている」に期待されるのは「成功」 
- 「タイトルまたは著者名で検索する」に対し「一致するものがない」に期待されるのは「空の結果」