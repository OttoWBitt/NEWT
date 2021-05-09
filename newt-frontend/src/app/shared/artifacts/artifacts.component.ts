import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { Artifact } from 'src/app/models/artifact.model';
import { ArtifactService } from 'src/app/services/views/artifact.service';

@Component({
  selector: 'newt-artifacts',
  templateUrl: './artifacts.component.html',
  styleUrls: ['./artifacts.component.scss']
})
export class ArtifactsComponent implements OnInit, AfterViewInit{

  fileToUpload: File = null;
  artifact: Artifact = new Artifact();
  artifacts: Artifact[];

  displayedColumns: string[] = ['name', 'subject', 'link', 'username', 'download'];
  dataSource = new MatTableDataSource<Artifact>();

  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(
    private artifactService: ArtifactService
  ) { }

  ngOnInit(): void {
    this.loadData()
  }

  loadData(){
    this.artifactService.getAllArtifacts().subscribe(response => {
      this.dataSource.data = response
      console.log(response)
    })
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
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