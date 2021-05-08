
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Constants } from 'src/app/util/constants';
import { Response } from 'src/app/models/response.model'
import { Artifact } from 'src/app/models/artifact.model';

@Injectable()
export class ArtifactService {

  constructor(private httpClient: HttpClient) { }
  
  newArtifact(artifact: Artifact) : Observable<any> {
    const formData: FormData = new FormData();
    formData.append('artifactFile', artifact.file);
    formData.append('artifactLink', artifact.link);
    formData.append('artifactName', artifact.name);
    formData.append('artifactDescription', artifact.description);
    formData.append('artifactUserId', artifact.user.id.toString());
    formData.append('artifactSubjectId', artifact.subject.id.toString());

    return this.httpClient.post<Response>(`${Constants.BASE_URL}artifact/new`, formData, {headers: this.getHeaders()}).pipe(
      map((response: Response) => {
        const resp: any = response;
        return resp;
      }));
  }

  private getHeaders(): HttpHeaders {
    const headers = new HttpHeaders();
    return headers;
  }

}