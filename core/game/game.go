package game

var (
	Marks = [2]string{"X", "O"}
)

func CheckResult(field *[3][3]string) string {
	var (
		isDraw = true
		isEnd  = false
		player string
	)

	for i := 0; i < 3; i++ {
		if field[i][0] != "" && field[i][0] == field[i][1] && field[i][1] == field[i][2] {
			player = field[i][0]
			isEnd = true
		}
		if field[0][i] != "" && field[0][i] == field[1][i] && field[1][i] == field[2][i] {
			player = field[0][i]
			isEnd = true
		}
	}
	if field[0][0] != "" && field[0][0] == field[1][1] && field[1][1] == field[2][2] {
		player = field[0][0]
		isEnd = true
	}
	if field[0][2] != "" && field[0][2] == field[1][1] && field[1][1] == field[2][0] {
		player = field[0][2]
		isEnd = true
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if field[i][j] == "" {
				isDraw = false
			}
		}
	}

	if isEnd {
		return player
	}

	if isDraw {
		return "draw"
	}

	return "ongoing"
}

func Validate(field *[3][3]string, x, y int) bool {
	if x < 0 || x > 2 {
		return false
	}
	if y < 0 || y > 2 {
		return false
	}
	return field[x][y] == ""
}
