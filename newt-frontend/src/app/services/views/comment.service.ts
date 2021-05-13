
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Commentary } from 'src/app/models/commentary.model';
import { Response } from 'src/app/models/response.model';
import { Constants } from 'src/app/util/constants';

@Injectable()
export class CommentService {

  constructor(private httpClient: HttpClient) { }
    
  getCommentsByArtifactId(artifactId: number) : Observable<Commentary[]> {
    return this.httpClient.get<any>(`${Constants.BASE_URL}comment/all/artifact/${artifactId}`, {headers: this.getHeaders()}).pipe(
      map((response: Response) => {
        const resp: Commentary[] = response.data;
        return resp;
      }));
  }
    
  saveComment(comment: Commentary) : Observable<Commentary> {
    return this.httpClient.post<any>(`${Constants.BASE_URL}comment/new`, comment,{headers: this.getHeaders()}).pipe(
      map((response: Response) => {
        const resp: Commentary = response.data;
        return resp;
      }));
  }

  private getHeaders(): HttpHeaders {
    const headers = new HttpHeaders();
    return headers;
  }

}