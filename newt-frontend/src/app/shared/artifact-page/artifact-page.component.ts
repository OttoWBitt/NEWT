import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { Artifact } from 'src/app/models/artifact.model';
import { Commentary } from 'src/app/models/commentary.model';
import { ArtifactService } from 'src/app/services/views/artifact.service';
import { CommentService } from 'src/app/services/views/comment.service';

@Component({
  selector: 'newt-artifact-page',
  templateUrl: './artifact-page.component.html',
  styleUrls: ['./artifact-page.component.scss']
})
export class ArtifactPageComponent implements OnInit {

  artifactId: number
  artifact: Artifact = new Artifact()
  comment: Commentary = new Commentary()
  comments: Commentary[]
  commentForm: FormGroup

  constructor(
    private formBuilder: FormBuilder,
    private artifactService: ArtifactService,
    private commentService: CommentService,
    private route: ActivatedRoute
    ) { }

  ngOnInit(): void {
    this.initForm()
    this.loadData()
  }

  initForm(){
    this.commentForm = this.formBuilder.group({
      comment: [null]
    });
  }

  loadData(){
    this.artifactId = parseInt(this.route.snapshot.paramMap.get('id'))
    this.artifactService.getArtifactById(this.artifactId).subscribe(response => {
      if (response) {
        this.artifact = response
      }
    })
    this.commentService.getCommentsByArtifactId(this.artifactId).subscribe(response => {
      if (response) {
        this.comments = response
      }
    })
  }

  addComment(){
    this.comment = this.commentForm.get('commment').value
    this.comment.user.id = this.getUserId()
    this.comment.artifact.id = this.artifactId
    this.commentService.saveComment(this.comment).subscribe(response => {
      if (response){
        this.loadData()
      }
    })
  }
  
  goToUrl(link: string){
    window.open(link, "_blank");
  }

  getUserId(){
    const sessionData = JSON.parse(localStorage.getItem('currentUser'));
    return sessionData.user.id
  }

}
