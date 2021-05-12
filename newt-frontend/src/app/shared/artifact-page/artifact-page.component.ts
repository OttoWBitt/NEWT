import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Artifact } from 'src/app/models/artifact.model';
import { ArtifactService } from 'src/app/services/views/artifact.service';

@Component({
  selector: 'newt-artifact-page',
  templateUrl: './artifact-page.component.html',
  styleUrls: ['./artifact-page.component.scss']
})
export class ArtifactPageComponent implements OnInit {

  artifactId: number
  artifact: Artifact = new Artifact()
  artifacts: Artifact[]

  constructor(
    private artifactService: ArtifactService,
    private route: ActivatedRoute
    ) { }

  ngOnInit(): void {
    this.loadData()
  }

  loadData(){
    this.artifactId = parseInt(this.route.snapshot.paramMap.get('id'))
    this.artifactService.getArtifactsById(this.artifactId).subscribe(response => {
      console.log(response)
      if (response) {
        this.artifact = response[0]

      }
    })
  }
  
  goToUrl(link: string){
    window.open(link, "_blank");
  }

}
