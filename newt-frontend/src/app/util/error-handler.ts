const errorsList = [
    {errorMessage: 'username not available', feedbackMessage: 'Nome de usuário não disponível'},
    {errorMessage: 'email not available', feedbackMessage: 'Email já cadastrado'},
]

export function generalExceptionTreatment(error){
    return errorsList.find(e => error.error.errors.includes(e.errorMessage)).feedbackMessage
}
