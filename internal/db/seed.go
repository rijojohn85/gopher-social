package db

import (
	"context"
	"log"
	"math/rand"
	"strconv"

	"github.com/rijojohn85/social/internal/store"
)

func Seed(store store.Storage) {
	ctx := context.Background()
	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating posts: ", err)
			return
		}
	}
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}
	}
	log.Println("Seeding complete")
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)

	for i := 0; i < count; i++ {
		users[i] = &store.User{
			Username: "user" + strconv.Itoa(i),
			Password: "password" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@gmail.com",
		}
	}
	return users
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = &store.Post{
			Title:   "title" + strconv.Itoa(i),
			Content: "content" + strconv.Itoa(i),
			Tags:    []string{"tag" + strconv.Itoa(i), "tag2" + strconv.Itoa(i)},
			UserID:  users[i%100].ID,
		}
	}
	return posts
}

func generateComments(count int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserId:  users[rand.Intn(len(users))].ID,
			Content: posts[rand.Intn(len(posts))].Content,
		}
	}
	return comments
}
