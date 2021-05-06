import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { LoginService } from 'src/app/services/views/login.service';

@Component({
  selector: 'newt-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.scss']
})
export class LandingComponent implements OnInit {

  loginForm: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private loginService: LoginService,
  ) { }

  ngOnInit(): void {
    this.initForm();
  }

  initForm() {
    this.loginForm = this.formBuilder.group({
      username: [null , Validators.required],
      password: [null , Validators.required],
    });
  }

  teste(){
    this.loginService.check().subscribe(r=>{
      console.log(r)
    })
  }

  signIn(){
    let formData = this.loginForm.getRawValue()
    this.loginService.login(formData.username, formData.password).subscribe(
      (response)=>{
        console.log('Login successful')
      },
      (error)=>{
        console.log(error)
      }
    )
    console.log(formData)
  }

  passwordRecovery(){
    console.log('User has forgotten his password')
  }
}
