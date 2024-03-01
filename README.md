# GoClinic Project  
## Project Description  
This project is a web-based appointment scheduling system that facilitates the process of scheduling appointments between patients and doctors. It provides an efficient and convenient way for patients to book appointments with their healthcare providers, reducing wait times and improving overall scheduling efficiency.
## Restaurants REST API  
&nbsp;POST /creation  
&nbsp;GET /registrations/:id  
&nbsp;PUT /registrations/:id  
&nbsp;DELETE /registrations/:id  

## DB Structure  
Table doctors {  
&nbsp;&nbsp;&nbsp;&nbsp;     id bigserial [primary key]  
&nbsp;&nbsp;&nbsp;&nbsp;     updated_at timestamp      created_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     updated_at timestamp  
&nbsp;&nbsp;&nbsp;&nbsp;     first_name text  
&nbsp;&nbsp;&nbsp;&nbsp;     last_name text  
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

## Team Members  
&nbsp;&nbsp;&nbsp;&nbsp;    Zhaxylykuly Aidar 22В030538  
&nbsp;&nbsp;&nbsp;&nbsp;    Zhaslan Bexultan  &nbsp;22B030355













