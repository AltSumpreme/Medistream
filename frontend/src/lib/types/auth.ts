export type JwtPayload = {
    id: string;
    role: 'PATIENT'| 'DOCTOR' | 'ADMIN';
    exp : number;
}