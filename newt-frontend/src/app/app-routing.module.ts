import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from './services/auth/auth.guard';
import { ArtifactPageComponent } from './shared/artifact-page/artifact-page.component';
import { ArtifactsComponent } from './shared/artifacts/artifacts.component';
import { SystemLayoutComponent } from './shared/system-layout/system-layout.component';
import { Views } from './util/views.enum';


const routes: Routes = [
  {
    path: '',
    component: SystemLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      {
        path: Views.artifacts.url,
        component: ArtifactsComponent
      },
      {
        path: Views.artifact.url,
        component: ArtifactPageComponent
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  providers: [AuthGuard],
  exports: [RouterModule]
})
export class AppRoutingModule { }
