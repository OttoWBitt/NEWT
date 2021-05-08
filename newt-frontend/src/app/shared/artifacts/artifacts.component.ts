import { Component, OnInit } from '@angular/core';
import { Artifact } from 'src/app/models/artifact.model';
import { ArtifactService } from 'src/app/services/views/artifact.service';

@Component({
  selector: 'newt-artifacts',
  templateUrl: './artifacts.component.html',
  styleUrls: ['./artifacts.component.scss']
})
export class ArtifactsComponent implements OnInit {

  fileToUpload: File = null;
  artifact: Artifact = new Artifact();

  constructor(
    private artifactService: ArtifactService
  ) { }

  ngOnInit(): void {
  }

  handleFileInput(files: FileList) {
    this.fileToUpload = files.item(0);
  }

  uploadFileToActivity() {
    this.artifact.user.id = 1
    this.artifact.name = 'Artefato Teste'
    this.artifact.subject.description = 'Geral'
    this.artifact.file = this.fileToUpload
    this.artifactService.newArtifact(this.artifact).subscribe(data => {
      // do something, if upload success
      }, error => {
        console.log(error);
      });
  }


}
