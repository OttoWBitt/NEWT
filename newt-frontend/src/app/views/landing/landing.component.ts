import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/views/auth.service';
import { Views } from 'src/app/util/views.enum';

@Component({
  selector: 'newt-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.scss']
})
export class LandingComponent implements OnInit {

  loginForm: FormGroup;
  recoverForm: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private router: Router,
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



  recoverPassword(){
    let formData = this.recoverForm.getRawValue()

    console.log(formData.email)
  }
}
