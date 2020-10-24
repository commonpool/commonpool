import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {ActivatedRoute, Router} from '@angular/router';
import {pluck} from 'rxjs/operators';

@Component({
  selector: 'app-resource-inquiry',
  templateUrl: './resource-inquiry.component.html',
  styleUrls: ['./resource-inquiry.component.css']
})
export class ResourceInquiryComponent implements OnInit {

  constructor(private backend: BackendService, private route: ActivatedRoute, private router: Router) {
  }

  pending = false;
  error: any = undefined;
  content = '';
  resource$ = this.route.params.pipe(pluck('id'));

  ngOnInit(): void {
  }

  submit(content: string, resource: string) {
    this.pending = true;
    this.error = undefined;
    this.backend.inquireAboutResource(resource, content).subscribe((res) => {
        this.pending = false;
        this.router.navigateByUrl('/messages').then(() => {
          console.log('OK!');
        });
      }, err => {
        console.error(err);
        this.pending = false;
        this.error = err;
      }
    );
  }
}
