package main

func vec_add(a, b vec) vec {
	return vec{a.X + b.X, a.Y + b.Y}
}

func vec_sub(a, b vec) vec {
	return vec{a.X - b.X, a.Y - b.Y}
}

func vec_mul(a, b vec) vec {
	return vec{a.X * b.X, a.Y * b.Y}
}

func vec_div(a, b vec) vec {
	return vec{a.X / b.X, a.Y / b.Y}
}
