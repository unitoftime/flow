package main

// func main() {
// 	fmt.Println("Starting")
// 	glitch.Run(runGame)
// }

// func runGame() {
// 	win, err := glitch.NewWindow(1920, 1080, "Dungeon Generation", glitch.WindowConfig{
// 		Vsync: false,
// 		Samples: 4,
// 	})
// 	if err != nil { panic(err) }

// 	shader, err := glitch.NewShader(shaders.SpriteShader)
// 	if err != nil { panic(err) }

// 	pass := glitch.NewRenderPass(shader)
// 	camera := glitch.NewCameraOrtho()

// 	rng := rand.New(rand.NewSource(99))
// 	dag, roomPos := pgen.GenerateRandomGridWalkDag2(rng, 100, 10)
// 	startLabel := "0_0"
// 	rooms := make(map[string]pgen.RoomPlacement)

// 	// gridSize := 10 // Min Room size
// 	for _, k := range dag.Nodes {
// 		static := (k == startLabel)

// 		rect := tile.R(-(rng.Intn(10)+5), -(rng.Intn(10)+5), rng.Intn(10)+5, rng.Intn(10)+5)

// 		// pos, ok := roomPos[k]
// 		// if !ok { panic("AAA") }
// 		// rect = rect.WithCenter(tile.TilePosition{pos.X * gridSize, pos.Y * gridSize})
// 		fmt.Println(rect.Center())
// 		rooms[k] = pgen.RoomPlacement{
// 			Rect: rect,
// 			GoalGap: math.Max(float64(rect.W()/2), float64(rect.H()/2)),
// 			// GoalGap: 10,
// 			Static: static,
// 			Placed: true,
// 			Repel: float64(rect.W() * rect.H()),// 1,
// 			Attract: 1.0,
// 		}
// 	}

// 	layout := pgen.NewGridLayout(rng, dag, rooms, startLabel, 5)
// 	layout.LayoutGrid(roomPos)

// 	zoom := 2.0

// 	phase := 0

// 	geom := glitch.NewGeomDraw()
// 	for !win.JustPressed(glitch.KeyEscape) {
// 		_, scrollY := win.MouseScroll()
// 		if scrollY < 0 {
// 			zoom  = zoom / 2.0
// 		} else if scrollY > 0 {
// 			zoom  = zoom * 2.0
// 		}

// 		pass.Clear()

// 		if phase == 0 {
// 			done := layout.IterateGravity()
// 			if done { phase = 1 }
// 		} else if phase == 1 {
// 			layout.IterateTowardsParent()
// 		}
// 		// layout.Iterate()
// 		// layout.Expand()

// 		for label, p := range rooms {
// 			if !p.Placed { continue } // Skip if not placed

// 			if pgen.HasRectIntersections(rooms, label) {
// 				geom.SetColor(glitch.RGBA{1, 0, 0, 1})
// 			} else {
// 				geom.SetColor(glitch.Greyscale(0.4))
// 			}

// 			rect := glitch.R(
// 				float64(p.Rect.Min.X),
// 				float64(p.Rect.Min.Y),
// 				float64(p.Rect.Max.X),
// 				float64(p.Rect.Max.Y),
// 			)
// 			mesh := geom.FillRect(rect)
// 			mesh.Draw(pass, glitch.Mat4Ident)
// 		}

// 		for node, edges := range dag.Edges {
// 			placement, ok := rooms[node]
// 			if !ok { panic("AAA") }
// 			if !placement.Placed { continue } // Skip if not placed

// 			srcPos := glitch.Vec3{
// 				float64(placement.Rect.Center().X),
// 				float64(placement.Rect.Center().Y),
// 				0,
// 			}

// 			for _, e := range edges {
// 				if pgen.HasEdgeIntersections(dag, rooms, node, e) {
// 					geom.SetColor(glitch.RGBA{1, 1, 0, 1})
// 				} else {
// 					geom.SetColor(glitch.Greyscale(0.9))
// 				}

// 				p2, ok := rooms[e]
// 				if !ok { panic("AAA") }
// 				if !p2.Placed { continue } // Skip if not placed

// 				dstPos := glitch.Vec3{
// 					float64(p2.Rect.Center().X),
// 					float64(p2.Rect.Center().Y),
// 					0,
// 				}
// 				mesh := glitch.NewMesh()
// 				geom.LineStrip(mesh, []glitch.Vec3{srcPos, dstPos}, 1)
// 				mesh.Draw(pass, glitch.Mat4Ident)
// 			}
// 		}

// 		glitch.Clear(win, glitch.RGBA{0, 0, 0, 1.0})
// 		camera.SetOrtho2D(win.Bounds())
// 		center := win.Bounds().Center()

// 		camera.SetView2D(-center[0], -center[1], zoom, zoom)
// 		pass.SetUniform("projection", camera.Projection)
// 		pass.SetUniform("view", camera.View)
// 		pass.Draw(win)
// 		win.Update()
// 	}
// }

// // func runGame() {
// // 	win, err := glitch.NewWindow(1920, 1080, "Dungeon Generation", glitch.WindowConfig{
// // 		Vsync: true,
// // 		Samples: 4,
// // 	})
// // 	if err != nil { panic(err) }

// // 	shader, err := glitch.NewShader(shaders.SpriteShader)
// // 	if err != nil { panic(err) }

// // 	pass := glitch.NewRenderPass(shader)
// // 	camera := glitch.NewCameraOrtho()

// // 	// dag := pgen.NewRoomDag()
// // 	// dag.AddNode("spawn")
// // 	// dag.AddNode("normal1")
// // 	// dag.AddNode("normal2")
// // 	// dag.AddNode("normal3")
// // 	// dag.AddNode("normal4")
// // 	// dag.AddNode("normal5")
// // 	// dag.AddNode("normal6")
// // 	// dag.AddNode("normal7")
// // 	// dag.AddNode("normal8")
// // 	// dag.AddNode("normal9")
// // 	// dag.AddNode("normal10")
// // 	// dag.AddNode("normal11")

// // 	// dag.AddEdge("spawn", "normal1")
// // 	// dag.AddEdge("normal1", "normal2")
// // 	// dag.AddEdge("normal2", "normal3")
// // 	// dag.AddEdge("normal3", "normal4")
// // 	// dag.AddEdge("normal3", "normal5")
// // 	// dag.AddEdge("normal5", "normal6")
// // 	// dag.AddEdge("normal6", "normal7")
// // 	// dag.AddEdge("normal5", "normal8")
// // 	// dag.AddEdge("normal8", "normal9")
// // 	// dag.AddEdge("normal9", "normal10")
// // 	// dag.AddEdge("normal10", "normal11")
// // 	// dag.AddEdge("normal8", "normal11")
// // 	// dag.AddEdge("normal11", "boss")

// // 	// startLabel := "spawn"
// // 	dag := pgen.GenerateRandomGridWalkDag(100, 3)
// // 	startLabel := "0_0"
// // 	rooms := make(map[string]pgen.RoomPlacement)

// // 	for k := range dag.Nodes {
// // 		static := (k == startLabel)

// // 		// rect := tile.R(-rand.Intn(50)+50, -rand.Intn(50)+50, rand.Intn(50)+50, rand.Intn(50)+50)
// // 		rect := tile.R(-rand.Intn(10)+10, -rand.Intn(10)+10, rand.Intn(10)+10, rand.Intn(10)+10)
// // 		rect = rect.WithCenter(tile.TilePosition{0, 0})
// // 		fmt.Println(rect.Center())
// // 		rooms[k] = pgen.RoomPlacement{
// // 			Rect: rect,
// // 			GoalGap: math.Max(float64(rect.W()/2), float64(rect.H()/2)),
// // 			// GoalGap: 10,
// // 			Static: static,
// // 			Placed: static,
// // 			Repel: float64(rect.W() * rect.H()),// 1,
// // 			Attract: 1.0,
// // 		}
// // 	}

// // 	rng := rand.New(rand.NewSource(99))
// // 	// pgen.PlaceDepthFirst(rng, dag, rooms, startLabel, 100)
// // 	// layout := pgen.NewDepthFirstLayout(rng, dag, rooms, startLabel, 100)
// // 	// layout := pgen.NewBreadthFirstLayout(rng, dag, rooms, startLabel, 100)

// // 	layout := pgen.NewHeirarchicalLayout(rng, dag, rooms, startLabel)
// // 	layout.GroupLayers(startLabel)

// // 	repel := 90000.0
// // 	grav := 0.0
// // 	relaxer := pgen.NewForceBasedRelaxer(rng, dag, rooms, repel, grav)

// // 	zoom := 1.0
// // 	phase := 0
// // 	// phaseCount := 0

// // 	geom := glitch.NewGeomDraw()
// // 	for !win.JustPressed(glitch.KeyEscape) {
// // 		_, scrollY := win.MouseScroll()
// // 		if scrollY < 0 {
// // 			zoom  = zoom / 2.0
// // 		} else if scrollY > 0 {
// // 			zoom  = zoom * 2.0
// // 		}

// // 		pass.Clear()

// // 		if phase == 0 {
// // 			cont := layout.Iterate()

// // 			if !cont {
// // 				cut := layout.CutCrossed()
// // 				if !cut {
// // 					phase++
// // 				}

// // 				// repel := 100000.0
// // 				// gravity := 0.0
// // 				// pgen.Untangle(dag, rooms, repel, gravity, 100)
// // 			}

// // 			// if !cont {
// // 			// 	stillCrossed := false
// // 			// 	for label := range rooms {
// // 			// 		if pgen.NodeHasEdgeIntersections(dag, rooms, label) {
// // 			// 			stillCrossed = true
// // 			// 		}
// // 			// 	}
// // 			// 	if stillCrossed {
// // 			// 		layout.Reset()
// // 			// 		layout = pgen.NewDepthFirstLayout(rng, dag, rooms, startLabel, 100)
// // 			// 		// layout = pgen.NewBreadthFirstLayout(rng, dag, rooms, startLabel, 100)
// // 			// 	} else {
// // 			// 		phase++
// // 			// 	}
// // 			// }

// // 			// cont := layout.Iterate()
// // 			// if !cont {
// // 			// 	for k, p := range rooms {
// // 			// 		if !p.Static {
// // 			// 			p.Placed = false
// // 			// 			rooms[k] = p
// // 			// 		}
// // 			// 	}

// // 			// 	layout = pgen.NewDepthFirstLayout(rng, dag, rooms, startLabel, 1)
// // 			// }

// // 			// noIntersections := true
// // 			// for node, edges := range dag.Edges {
// // 			// 	for _, e := range edges {
// // 			// 		if pgen.HasEdgeIntersections(dag, rooms, node, e) {
// // 			// 			r := rooms[node]
// // 			// 			r.Attract = 1.5
// // 			// 			r.Repel += 50
// // 			// 			rooms[node] = r

// // 			// 			noIntersections = false
// // 			// 		} else {
// // 			// 			r := rooms[node]
// // 			// 			// r.Attract = 1.0
// // 			// 			// r.Repel = float64(r.Rect.W() * r.Rect.H())
// // 			// 			rooms[node] = r
// // 			// 		}
// // 			// 	}
// // 			// }

// // 			// // pgen.PSLDStep(rng, dag, rooms)
// // 			// // time.Sleep(100 * time.Millisecond)

// // 			// repel := 100000.0
// // 			// gravity := 0.0
// // 			// pgen.Wiggle(dag, rooms, repel, gravity, 1)

// // 			// // if stable && noIntersections {
// // 			// if noIntersections {
// // 			// 	for label, r := range rooms {
// // 			// 		r.Attract = 1.0
// // 			// 		r.Repel = float64(r.Rect.W() * r.Rect.H())
// // 			// 		rooms[label] = r
// // 			// 	}

// // 			// 	phase++
// // 			// }

// // 			// phaseCount++
// // 			// if phaseCount > 100 {
// // 			// 	phaseCount = 0

// // 			// 	for label, r := range rooms {
// // 			// 		r.Attract = 1.0
// // 			// 		r.Repel = float64(r.Rect.W() * r.Rect.H())
// // 			// 		rooms[label] = r
// // 			// 	}
// // 			// 	repel := 10000.0
// // 			// 	gravity := 0.2
// // 			// 	pgen.Wiggle(dag, rooms, repel, gravity, 100)
// // 			// }
// // 		} else {
// // 			relaxer.Iterate()
// // 			// // noIntersections := true
// // 			// for node, edges := range dag.Edges {
// // 			// 	for _, e := range edges {
// // 			// 		if pgen.HasEdgeIntersections(dag, rooms, node, e) {
// // 			// 			r := rooms[node]
// // 			// 			r.Attract = 1.0
// // 			// 			r.Repel += 1
// // 			// 			rooms[node] = r

// // 			// 			// noIntersections = false
// // 			// 		} else {
// // 			// 			r := rooms[node]
// // 			// 			// r.Attract = 1.0
// // 			// 			// r.Repel = float64(r.Rect.W() * r.Rect.H())
// // 			// 			rooms[node] = r
// // 			// 		}
// // 			// 	}
// // 			// }

// // 			// // for label, p := range rooms {
// // 			// // 	if pgen.HasRectIntersections(rooms, label) {
// // 			// // 		// p.GoalGap++
// // 			// // 		p.Repel += 1
// // 			// // 		rooms[label] = p
// // 			// // 		break
// // 			// // 	} else {
// // 			// // 		// p.GoalGap--
// // 			// // 		// rooms[label] = p
// // 			// // 	}
// // 			// // }

// // 			// repel := 50000.0
// // 			// gravity := 0.1
// // 			// pgen.Wiggle(dag, rooms, repel, gravity, 1)
// // 		}

// // 		// if phase == 0 {
// // 		// 	// // noIntersections := true
// // 		// 	// for node, edges := range dag.Edges {
// // 		// 	// 	for _, e := range edges {
// // 		// 	// 		if pgen.HasEdgeIntersections(dag, rooms, node, e) {
// // 		// 	// 			r := rooms[node]
// // 		// 	// 			// r.Attract += 0.01
// // 		// 	// 			r.Attract = 2.0
// // 		// 	// 			rooms[node] = r

// // 		// 	// 			// noIntersections = false
// // 		// 	// 		} else {
// // 		// 	// 			r := rooms[node]
// // 		// 	// 			r.Attract = 1.0
// // 		// 	// 			rooms[node] = r
// // 		// 	// 		}
// // 		// 	// 	}
// // 		// 	// }
// // 		// 	repel := 10000.0
// // 		// 	gravity := 2.0
// // 		// 	stable := pgen.Wiggle(dag, rooms, repel, gravity, 1)
// // 		// 	// if stable || noIntersections {
// // 		// 	phaseCount++
// // 		// 	if stable || phaseCount > 500 {
// // 		// 		for label, p := range rooms {
// // 		// 			p.Repel = 50.0
// // 		// 			rooms[label] = p
// // 		// 		}

// // 		// 		phase++
// // 		// 	}
// // 		// } else if phase == 1 {
// // 		// 	noIntersections := true
// // 		// 	for label, p := range rooms {
// // 		// 		if pgen.HasRectIntersections(rooms, label) {
// // 		// 			p.GoalGap++
// // 		// 			// p.Repel += 1
// // 		// 			rooms[label] = p
// // 		// 			noIntersections = false
// // 		// 		}
// // 		// 	}
// // 		// 	repel := 600.0
// // 		// 	gravity := 1.0
// // 		// 	stable := pgen.Wiggle(dag, rooms, repel, gravity, 100)
// // 		// 	if stable && noIntersections {
// // 		// 		phase++
// // 		// 	}
// // 		// }

// // 		for label, p := range rooms {
// // 			if !p.Placed { continue } // Skip if not placed

// // 			if pgen.HasRectIntersections(rooms, label) {
// // 				geom.SetColor(glitch.RGBA{1, 0, 0, 1})
// // 			} else {
// // 				geom.SetColor(glitch.Greyscale(0.4))
// // 			}

// // 			rect := glitch.R(
// // 				float64(p.Rect.Min.X),
// // 				float64(p.Rect.Min.Y),
// // 				float64(p.Rect.Max.X),
// // 				float64(p.Rect.Max.Y),
// // 			)
// // 			mesh := geom.FillRect(rect)
// // 			mesh.Draw(pass, glitch.Mat4Ident)
// // 		}

// // 		for node, edges := range dag.Edges {
// // 			placement, ok := rooms[node]
// // 			if !ok { panic("AAA") }
// // 			if !placement.Placed { continue } // Skip if not placed

// // 			srcPos := glitch.Vec3{
// // 				float64(placement.Rect.Center().X),
// // 				float64(placement.Rect.Center().Y),
// // 				0,
// // 			}

// // 			for _, e := range edges {
// // 				if pgen.HasEdgeIntersections(dag, rooms, node, e) {
// // 					geom.SetColor(glitch.RGBA{1, 1, 0, 1})
// // 				} else {
// // 					geom.SetColor(glitch.Greyscale(0.9))
// // 				}

// // 				p2, ok := rooms[e]
// // 				if !ok { panic("AAA") }
// // 				if !p2.Placed { continue } // Skip if not placed

// // 				dstPos := glitch.Vec3{
// // 					float64(p2.Rect.Center().X),
// // 					float64(p2.Rect.Center().Y),
// // 					0,
// // 				}
// // 				mesh := glitch.NewMesh()
// // 				geom.LineStrip(mesh, []glitch.Vec3{srcPos, dstPos}, 1)
// // 				mesh.Draw(pass, glitch.Mat4Ident)
// // 			}
// // 		}

// // 		glitch.Clear(win, glitch.RGBA{0, 0, 0, 1.0})
// // 		camera.SetOrtho2D(win.Bounds())
// // 		center := win.Bounds().Center()

// // 		camera.SetView2D(-center[0], -center[1], zoom, zoom)
// // 		pass.SetUniform("projection", camera.Projection)
// // 		pass.SetUniform("view", camera.View)
// // 		pass.Draw(win)
// // 		win.Update()
// // 	}
// // }
