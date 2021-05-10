import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { Artifact } from 'src/app/models/artifact.model';
import { Subject } from 'src/app/models/subject.model';
import { ArtifactService } from 'src/app/services/views/artifact.service';
import { SubjectService } from 'src/app/services/views/subject.service';
import { ArtifactDialogComponent } from '../artifact-dialog/artifact-dialog.component';

@Component({
  selector: 'newt-artifacts',
  templateUrl: './artifacts.component.html',
  styleUrls: ['./artifacts.component.scss']
})
export class ArtifactsComponent implements OnInit, AfterViewInit{

  artifacts: Artifact[];
  subjects: Subject[];
  dialogRef: any;

  displayedColumns: string[] = ['name', 'description','subject', 'username', 'link', 'download'];
  dataSource = new MatTableDataSource<Artifact>();

  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(
    private artifactService: ArtifactService,
    private subjectService: SubjectService,
    private dialog: MatDialog
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

  addArtifact(){
    this.openDialog(ArtifactDialogComponent, null)
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

  openDialog(component, data): void {
    this.dialogRef = this.dialog.open(component, {
      maxWidth: '1200px',
      maxHeight: '800px',
      data: data
    });
  }
}