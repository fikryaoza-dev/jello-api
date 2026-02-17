package mapper

type Mapper[D any, M any] interface {
	ToDomain(M) D
	ToModel(D) M
}

type RowsScanner interface {
	Next() bool
	ScanDoc(dest interface{}) error
	Close() error
}

func ScanAndMap[M any, D any](
	rows RowsScanner,
	mp Mapper[D, M],
) ([]D, error) {

	var result []D

	for rows.Next() {
		var modelData M

		if err := rows.ScanDoc(&modelData); err != nil {
			return nil, err
		}

		result = append(result, mp.ToDomain(modelData))
	}

	return result, nil
}
