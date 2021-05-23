import { Artifact } from "./artifact.model";
import { User } from "./user.model";

export class Commentary {
    id: number;
    date: string;
    user: User = new User()
    artifact: Artifact = new Artifact()
    comment: string
}