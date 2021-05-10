import { Component, Input, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { Artifact } from 'src/app/models/artifact.model';
import { Subject } from 'src/app/models/subject.model';
import { ArtifactService } from 'src/app/services/views/artifact.service';
import { SubjectService } from 'src/app/services/views/subject.service';

@Component({
  selector: 'newt-artifact-dialog',
  templateUrl: './artifact-dialog.component.html',
  styleUrls: ['./artifact-dialog.component.scss']
})
export class ArtifactDialogComponent implements OnInit {

  fileToUpload: File = null;
  artifact: Artifact = new Artifact();
  subjects: Subject[];
  artifactForm: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private artifactService: ArtifactService,
    private subjectService: SubjectService,
    private dialogRef: MatDialogRef<ArtifactDialogComponent>,
  ) { }

  ngOnInit(): void {
    this.initForm();
    this.subjectService.getAllSubjects().subscribe(response => {
      if (response){
        this.subjects = response
      }
    })
  }
 
  initForm() {
    this.artifactForm = this.formBuilder.group({
      name: [null , Validators.required],
      description: [null , Validators.required],
      subject: [null , Validators.required],
      link: [null],
      file: [null]
    });
  }

  handleFileInput(files: FileList) {
    this.fileToUpload = files.item(0);
  }

  onSubmit() {
    if (this.artifactForm.valid){
      this.setArtifactData()
      this.artifactService.newArtifact(this.artifact).subscribe(response => {
        if (response){
          this.dialogRef.close()
        }
        }, error => {
          console.log(error);
      });
    }
  }

  setArtifactData(){
      let formData = this.artifactForm.getRawValue()
      this.artifact.user.id = this.getUserId()
      this.artifact.name = formData.name
      this.artifact.description = formData.description
      this.artifact.name = formData.name
      this.artifact.subject.id = formData.subject
      this.artifact.link = formData.link
      this.artifact.file = this.fileToUpload
  }

  getUserId(){
    const sessionData = JSON.parse(localStorage.getItem('currentUser'));
    return sessionData.id
  }

  close(){
    this.dialogRef.close()
  }
}
