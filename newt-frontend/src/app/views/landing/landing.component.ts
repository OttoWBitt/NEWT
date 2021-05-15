import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatExpansionPanel } from '@angular/material/expansion';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/views/auth.service';
import { MessagesEnum } from 'src/app/util/messages.enum';
import { Views } from 'src/app/util/views.enum';

@Component({
  selector: 'newt-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.scss']
})
export class LandingComponent implements OnInit {

  @ViewChild(MatExpansionPanel) expansionPanel: MatExpansionPanel ;

  loginForm: FormGroup;
  recoverForm: FormGroup;
  panelOpenState: boolean;

  constructor(
    private formBuilder: FormBuilder,
    private router: Router,
    private snackbar: MatSnackBar,
    private authService: AuthService,
  ) { }

  ngOnInit(): void {
    this.initForm();
  }

  initForm() {
    this.loginForm = this.formBuilder.group({
      username: [null , Validators.required],
      password: [null , Validators.required],
    });
    this.recoverForm = this.formBuilder.group({
      email: [null , [Validators.required, Validators.pattern("^[a-z0-9._%+-]+@[a-z0-9.-]+\\.[a-z]{2,4}$")]]
    });
  }

  signIn(){
    let formData = this.loginForm.getRawValue()
    this.authService.login(formData.username, formData.password).subscribe(
      (response)=>{
        this.router.navigate([Views.artifacts.url])
      },
      (error)=>{
        console.log(error)
      }
    )
  }

  goToSignUp(){
    this.router.navigate([Views.signUp.url])
  }
  
  recoverPassword(){
    let formData = this.recoverForm.getRawValue()
    if (this.recoverForm.valid){
      this.authService.recoverPassword(formData.email).subscribe(response => {
        if(response){
          this.snackbar.open(MessagesEnum.RecoveryEmailSent);
          this.expansionPanel.close()
        }
      },
      error => {

      }
      )
    }
  }
}
