import { Subject } from "./subject.model";
import { User } from "./user.model";

export class Artifact {
    id: number;
    name: string;
    description: string;
    user: User = new User();
    subject: Subject = new Subject();
    file: File
    link: string;
    downloadLink: string;
}
