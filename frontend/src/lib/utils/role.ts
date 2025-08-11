export function hasRole(user:{role:string}, ...roles:string[]){
    return roles.includes(user.role);
}