package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"unicode"

	pb "github.com/bruceoaudo/userService/gen/user"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	pb.UnimplementedUserServiceServer
	db *mongo.Client
}

type User struct {
	FullName     string    `bson:"full_name"`
	UserName     string    `bson:"user_name"`
	EmailAddress string    `bson:"email"`
	PhoneNumber  string    `bson:"phone"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

// LoginUser remains exactly the same
func (s *userService) LoginUser(ctx context.Context, req *pb.LoginMessageRequest) (*pb.LoginMessageResponse, error) {
	
	// 1. Find user by email
	collection := s.db.Database("userdb").Collection("users")
	var user User
	err := collection.FindOne(ctx, bson.M{"email": strings.ToLower(req.GetEmail())}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Invalid credentials")
		}
		log.Printf("Database error: %v", err)
		return nil, status.Error(codes.Internal, "login failed")
	}
	
	
	return &pb.LoginMessageResponse{
		Email:    user.EmailAddress,
		UserName: user.UserName,
		Password: user.PasswordHash,
	}, nil
}

// RegisterUser function
func (s *userService) RegisterUser(ctx context.Context, req *pb.RegisterMessageRequest) (*pb.RegisterMessageResponse, error) {
	// 1. Validate input
	if err := validateRegistration(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	collection := s.db.Database("userdb").Collection("users")

	// 2. Check for existing user
	existingFilter := bson.M{
		"$or": []bson.M{
			{"email": strings.ToLower(strings.TrimSpace(req.GetEmailAddress()))},
			{"user_name": strings.ToLower(strings.TrimSpace(req.GetUserName()))},
			{"phone": normalizePhoneNumber(req.GetPhoneNumber())},
		},
	}

	var existingUser User
	err := collection.FindOne(ctx, existingFilter).Decode(&existingUser)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email, username or phone already exists")
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Database error: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// 3. Create user document
	user := User{
		FullName:     strings.TrimSpace(req.GetFullName()),
		UserName:     strings.ToLower(strings.TrimSpace(req.GetUserName())),
		EmailAddress: strings.ToLower(strings.TrimSpace(req.GetEmailAddress())),
		PhoneNumber:  normalizePhoneNumber(req.GetPhoneNumber()),
		PasswordHash: req.GetPassword(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(codes.AlreadyExists, "user with these details already exists")
		}
		log.Printf("Failed to create user: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.RegisterMessageResponse{
		UserName: user.UserName,
		Message:  "Registered successfully",
		Success:  true,
	}, nil
}

// Initialize MongoDB connection
func NewUserService(mongoURI string) (*userService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Create unique indexes
	db := client.Database("userdb")
	collection := db.Collection("users")

	_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{primitive.E{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{primitive.E{Key: "user_name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{primitive.E{Key: "phone", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, err
	}

	return &userService{db: client}, nil
}


func validateRegistration(req *pb.RegisterMessageRequest) error {
	if strings.TrimSpace(req.GetFullName()) == "" {
		return errors.New("full name is required")
	}

	username := strings.TrimSpace(req.GetUserName())
	if len(username) < 4 {
		return errors.New("username must be at least 4 characters")
	}
	if !isAlphanumeric(username) {
		return errors.New("username can only contain letters and numbers")
	}

	email := strings.TrimSpace(req.GetEmailAddress())
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("invalid email format")
	}

	phone := normalizePhoneNumber(req.GetPhoneNumber())
	if len(phone) != 12 || !strings.HasPrefix(phone, "254") {
		return errors.New("phone must be in 254XXXXXXXXX format (12 digits)")
	}

	return nil
}

func normalizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")

	switch {
	case strings.HasPrefix(phone, "0") && len(phone) == 10:
		return "254" + phone[1:]
	case strings.HasPrefix(phone, "7") && len(phone) == 9:
		return "254" + phone
	default:
		return phone
	}
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set in .env file")
	}

	userSvc, err := NewUserService(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userSvc)

	log.Printf("gRPC server listening on port: %s", "50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
