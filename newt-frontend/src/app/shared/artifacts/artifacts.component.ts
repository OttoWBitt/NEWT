import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { Artifact } from 'src/app/models/artifact.model';
import { Subject } from 'src/app/models/subject.model';
import { ArtifactService } from 'src/app/services/views/artifact.service';
import { SubjectService } from 'src/app/services/views/subject.service';

@Component({
  selector: 'newt-artifacts',
  templateUrl: './artifacts.component.html',
  styleUrls: ['./artifacts.component.scss']
})
export class ArtifactsComponent implements OnInit, AfterViewInit{

  fileToUpload: File = null;
  artifact: Artifact = new Artifact();
  artifacts: Artifact[];
  subjects: Subject[];

  displayedColumns: string[] = ['name', 'subject', 'username', 'link', 'download'];
  dataSource = new MatTableDataSource<Artifact>();

  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(
    private artifactService: ArtifactService,
    private subjectService: SubjectService
  ) { }

  ngOnInit(): void {
    this.loadData()
  }

  loadData(){
    this.artifactService.getAllArtifacts().subscribe(response => {
      if (response){
        this.artifacts = response
        this.dataSource.data = this.artifacts
      }
    })
    this.subjectService.getAllSubjects().subscribe(response => {
      if (response){
        this.subjects = response
      }
    })
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();
  }

  filterSubject(subjectIdList: string[]){
    if (subjectIdList.length > 0){
      this.dataSource.data = this.artifacts.filter(artifact => subjectIdList.includes(artifact.subject.id.toString()))
    } else {
      this.dataSource.data = this.artifacts
    }
  }

  goToUrl(link: string){
    window.open(link, "_blank");
  }

  handleFileInput(files: FileList) {
    this.fileToUpload = files.item(0);
  }

  uploadFileToActivity() {
    this.artifact.user.id = 1
    this.artifact.name = 'Artefato Teste'
    this.artifact.subject.name = 'Geral'
    this.artifact.file = this.fileToUpload
    this.artifactService.newArtifact(this.artifact).subscribe(data => {
      // do something, if upload success
      }, error => {
        console.log(error);
      });
  }
}