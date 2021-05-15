import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { ArtifactsComponent } from './shared/artifacts/artifacts.component';
import { SystemLayoutComponent } from './shared/system-layout/system-layout.component';
import { LandingRoutingModule } from './views/landing/landing-routing.module';
import { LandingComponent } from './views/landing/landing.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatDialogModule } from '@angular/material/dialog';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { AuthService } from './services/views/auth.service';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { AuthInterceptor } from './services/auth/auth.interceptor';
import { ArtifactService } from './services/views/artifact.service';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorIntl, MatPaginatorModule } from '@angular/material/paginator';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { SubjectService } from './services/views/subject.service';
import { MatToolbarModule } from '@angular/material/toolbar';
import { HeaderComponent } from './shared/header/header.component';
import { ArtifactDialogComponent } from './shared/artifact-dialog/artifact-dialog.component';
import { MaterialFileInputModule } from 'ngx-material-file-input';
import { getPtPaginatorIntl } from './util/pt-paginator-intl';
import { ArtifactPageComponent } from './shared/artifact-page/artifact-page.component';
import { CommentService } from './services/views/comment.service';
import { MatSnackBarModule, MAT_SNACK_BAR_DEFAULT_OPTIONS } from '@angular/material/snack-bar';
import { MatExpansionModule } from '@angular/material/expansion';

@NgModule({
  declarations: [
    AppComponent,
    LandingComponent,
    SystemLayoutComponent,
    ArtifactsComponent,
    HeaderComponent,
    ArtifactDialogComponent,
    ArtifactPageComponent
  ],
  imports: [
    BrowserModule,
    LandingRoutingModule,
    AppRoutingModule,
    FlexLayoutModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatDividerModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    MatTableModule,
    MatPaginatorModule,
    MatCardModule,
    MatIconModule,
    MatSelectModule,
    MatToolbarModule,
    MatDialogModule,
    MaterialFileInputModule,
    MatSnackBarModule,
    MatExpansionModule,
    HttpClientModule,
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true },
    { provide: MatPaginatorIntl, useValue: getPtPaginatorIntl()},
    { provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: {duration: 5000}},
    AuthService,
    ArtifactService,
    CommentService,
    SubjectService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
