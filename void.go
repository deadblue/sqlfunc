package sqlfunc

// Void implements [sql.Scanner] interface but discards the value.
type Void struct{}

func (v Void) Scan(src any) error {
	return nil
}
