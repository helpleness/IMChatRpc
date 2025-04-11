import { defineStore } from 'pinia';

export const useChatStore = defineStore('chat', {
    // State: Defines the initial state of the store
    state: () => ({
        isLoggedIn: '',  // Login status
        username: '',  // Example username
        userID: '',
        onlineusers: ['用户1','用户2'],
        allusers: ['用户1', '用户2', '用户3'],  // Example user list
        messages: [],  // Chat messages
        newMessage: '',  // New message content
        isVideoCalling: false,  // Video call status
        selectedUser: "",
        client: null,
    }),

    // Actions: Define methods that can modify the state
    actions: {
        // Login function
        login(userID) {
            this.isLoggedIn = true;
            this.userId = userID;
        },

        // Logout function
        logout() {
            this.isLoggedIn = false;
        },

        // Send message function
        sendMessage(message) {
            this.messages.push({ user: this.username, content: message });
        },
        receiveMessage(userID,message){
            this.messages.push({user: userID, content: message});

        },

        getusername(username){
            this.username=username;
        },

        selected_User(selectedUser){
            this.selectedUser=selectedUser;
        },

        client_connect(client){
            this.client = client;
        },

        loadMessagesForUser(user){

        },

        fetchFriends(){


        },

        fetchFriendRequests(){

        }
    },
});
