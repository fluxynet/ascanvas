package ascanvas

func TransformRectangle(canvas *Canvas, args TransformRectangleArgs) error {
	var (
		maxX, maxY int
		grid       = canvas.AsGrid()
	)

	if args.Width == 0 {
		maxX = args.TopLeft.X
	} else {
		maxX = args.TopLeft.X + args.Width - 1
	}

	if maxX >= canvas.Width {
		maxX = canvas.Width - 1
	}

	if args.Height == 0 {
		maxY = args.TopLeft.Y
	} else {
		maxY = args.TopLeft.Y + args.Height - 1
	}

	if maxY >= canvas.Height {
		maxY = canvas.Height - 1
	}

	if args.Fill != "" {
		for x := args.TopLeft.X; x <= maxX; x++ {
			for y := args.TopLeft.Y; y <= maxY; y++ {
				grid[y][x] = args.Fill
			}
		}
	}

	if args.Outline != "" {
		for x := args.TopLeft.X; x <= maxX; x++ {
			grid[args.TopLeft.Y][x] = args.Outline
			grid[maxY][x] = args.Outline
		}

		for y := args.TopLeft.Y; y <= maxY; y++ {
			grid[y][args.TopLeft.X] = args.Outline
			grid[y][maxX] = args.Outline

		}
	}

	canvas.FromGrid(grid)

	return nil
}

func TransformFloodfill(canvas *Canvas, args TransformFloodfillArgs) error {
	if !canvas.Contains(args.Start) {
		return ErrOutOfBounds
	}

	var (
		grid    = canvas.AsGrid()
		pattern = grid[args.Start.Y][args.Start.X]
	)
	transformFloodfill(canvas, args, grid, pattern)
	return nil
}

func transformFloodfill(canvas *Canvas, args TransformFloodfillArgs, grid [][]string, pattern string) {
	var (
		queue   = []Coordinates{args.Start}
		visited = make(map[string]interface{})
	)

	for len(queue) != 0 {
		var (
			p = queue[0]
			s = p.String()
		)
		queue = queue[1:]

		if _, ok := visited[s]; ok || !canvas.Contains(p) {
			continue
		}

		visited[s] = nil

		if grid[p.Y][p.X] == pattern {
			grid[p.Y][p.X] = args.Fill
		} else {
			continue
		}

		var next = []Coordinates{
			{X: p.X, Y: p.Y - 1}, // N
			{X: p.X, Y: p.Y + 1}, // S
			{X: p.X + 1, Y: p.Y}, // E
			{X: p.X - 1, Y: p.Y}, // W
		}

		for i := range next {
			var _, ok = visited[next[i].String()]

			switch {
			case !canvas.Contains(next[i]):
				continue
			case ok:
				continue
			case grid[next[i].Y][next[i].X] != pattern:
				continue
			default:
				queue = append(queue, next[i])
			}
		}
	}

	canvas.FromGrid(grid)
}
