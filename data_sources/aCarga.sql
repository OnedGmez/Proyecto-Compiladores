CREATE FUNCTION ContraseniaRandom 
()
RETURNS VARCHAR(10)
AS
BEGIN    

/*
ESTO ES UN COMENTARIO MULTIPLE CONN ERROR
*

    DECLARE @chars AS VARCHAR(52),
            @numbers AS VARCHAR(10),
            @Caracteres AS VARCHAR(62),        
            @contrasenia AS VARCHAR(62),
            @index AS INT,
            @cont AS INT

    SET @contrasenia = ''
    SET @Caracteres = ''    
    SET @chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
    SET @numbers = '0123456789'

    SET @Caracteres = @chars + @numbers

    SET @cont = 0
    WHILE @cont < 20
    BEGIN
        SET @index = ceiling( ( SELECT rnd FROM RandomNumero ) * (len(@Caracteres)))
        SET @contrasenia = @contrasenia + substring(@Caracteres, @index, 1)
        SET @cont = @cont + 1
    END    
    *
    ESTO ES UN COMENTARIO MULTIPLE CON ERROR
    */
        
    RETURN @contrasenia

END
GO