import { Component, OnInit } from '@angular/core';
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
  comments: Commentary[]

  constructor(
    private artifactService: ArtifactService,
    private commentService: CommentService,
    private route: ActivatedRoute
    ) { }

  ngOnInit(): void {
    this.loadData()
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
  
  goToUrl(link: string){
    window.open(link, "_blank");
  }

}
