import { Artifact } from "./artifact.model";
import { User } from "./user.model";

export class Commentary {
    id: number;
    date: string;
    user: User
    artifact: Artifact
    comment: string
}