package lorekeeper_test

import (
	"log"
	"testing"

	"github.com/trviph/lorekeeper"
)

const lorem = "Culpa sequi esse et et expedita aut qui quia. Error minus modi sunt beatae asperiores qui rem. Quia minima cumque laudantium sed rerum. Sunt delectus nesciunt dolor veniam soluta provident porro deserunt. Ullam illo beatae et quos unde maxime repellendus. Beatae itaque totam eum itaque velit et. Sit molestias dolore deserunt rerum amet. Molestiae rem provident minima autem nulla numquam. Illum voluptas ea nam suscipit. Corporis molestias necessitatibus dolore facilis. Nostrum cum nemo vero. Enim dolorem esse ad. Sed numquam odio eum ex. Praesentium incidunt quod perferendis sit est omnis sapiente. Sed rem itaque laboriosam minus eos. Sed fugiat dolores ut. Nam veniam nihil voluptatem accusamus molestias ducimus. Minima aut consequuntur dolores facere inventore libero tempore omnis. Suscipit et aut nostrum. Porro sapiente dignissimos nisi error. Et nulla vel molestiae veniam molestiae eum. Est similique sapiente aperiam voluptate cum occaecati et laboriosam. Praesentium cupiditate et laboriosam aperiam neque ut ut. Provident blanditiis autem pariatur autem animi et sint dicta."

func BenchmarkKeeperWrite(b *testing.B) {
	k, _ := lorekeeper.NewKeeper(
		lorekeeper.WithFolder("."),
		lorekeeper.WithMaxSize(100*lorekeeper.KB),
		lorekeeper.WithMaxFiles(3),
	)
	logger := log.New(k, "[Benchmark] ", log.LstdFlags|log.Lmsgprefix)
	b.Run(
		"Write to Keeper",
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logger.Println(lorem)
			}
		},
	)
}
