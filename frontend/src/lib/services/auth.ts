import api  from "$lib/config/api";


const loginUser = async (email : string , password : string) => {
    const response = await api.post("/api/login", {
        email,
        password,
    });

    return response.data;
}


const registerUser = async (firstName: string, lastName: string, email: string, password: string) => {
    const response = await api.post("/api/signup", {
        firstName,
        lastName,
        email,
        password,
    });

    return response.data;
}

export { loginUser, registerUser };