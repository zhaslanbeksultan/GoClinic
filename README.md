# GoClinic Project  
## Restaurants REST API  
&nbsp;POST /menus  
&nbsp;GET /menus/:id  
&nbsp;PUT /menus/:id  
&nbsp;DELETE /menus/:id  

## DB Structure  
Table doctors {  
&nbsp;&nbsp;&nbsp;&nbsp;    id bigserial [primary key]  
&nbsp;&nbsp;&nbsp;&nbsp;    updated_at timestamp      created_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     updated_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     first_name text  
&nbsp;&nbsp;&nbsp;&nbsp;    last_name text  
&nbsp;&nbsp;&nbsp;&nbsp;     speciality text  
&nbsp;&nbsp;&nbsp;&nbsp;     phone text  
&nbsp;    } 

Table patients {  
&nbsp;&nbsp;&nbsp;&nbsp;     id bigserial [primary key]  
&nbsp;&nbsp;&nbsp;&nbsp;     created_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     updated_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     first_name text  
&nbsp;&nbsp;&nbsp;&nbsp;     last_name text  
&nbsp;&nbsp;&nbsp;&nbsp;     phone text  
&nbsp;    }  

// many-to-many  
Table appointments {  
&nbsp;&nbsp;&nbsp;&nbsp;     id bigserial [primary key]  
&nbsp;&nbsp;&nbsp;&nbsp;     created_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     updated_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     doctor_id bigserial  
&nbsp;&nbsp;&nbsp;&nbsp;     patient_id bigserial  
&nbsp;    }  

Ref: appointments.doctor_id < doctors.id  
Ref: appointments.patient_id < patients.id  
